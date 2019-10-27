package comm

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"

	json "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

func GetCurrentPath() string {
	_, filename, _, ok := runtime.Caller(1)
	var cwdPath string
	if ok {
		cwdPath = path.Join(path.Dir(filename), "") // the the main function file directory
	} else {
		cwdPath = "./"
	}
	return cwdPath
}

const getIPURL = "http://ipinfo.io/json"

// IPStruct ip address info structure
type IPStruct struct {
	IP      string `json:"ip"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	Loc     string `json:"loc"`
	Org     string `json:"org"`
}

// GetIP getting external IP from some service
func GetIP() string {
	resp, err := http.Get(getIPURL)
	if err != nil {
		logrus.Error("Error during get external IP info from:" + getIPURL)
		return "127.0.0.1"
	}
	defer resp.Body.Close()

	var newIP IPStruct
	json.NewDecoder(resp.Body).Decode(&newIP)

	return newIP.IP
}

// GetIPInfo getting external IP full info from some service
func GetIPInfo() IPStruct {
	resp, err := http.Get(getIPURL)
	if err != nil {
		fmt.Println("Error during get external IP info from:" + getIPURL)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var newIP IPStruct
	json.NewDecoder(resp.Body).Decode(&newIP)

	return newIP
}
