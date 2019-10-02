package main

import (
	"fmt"

	"github.com/json-iterator/go"
)

func main() {
	// 完全兼容标准json库
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	u := User{
		Id:      12,
		Name:    "wendell",
		Age:     1,
		Address: "成都高新区",
	}
	data, err := json.Marshal(&u)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))

	u2 := &User{}
	err = json.Unmarshal(data, u2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v \n", u)
	fmt.Printf("%+v \n", u2)
}
