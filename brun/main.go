package main

import (
	"github.com/tietang/props/ini"

	// 调用app.go的init方法初始化各种资源启动器
	_ "github.com/solozyx/red-envelope"
	"github.com/solozyx/red-envelope/comm"
	"github.com/solozyx/red-envelope/infra"
)

func main() {
	// 获取程序运行文件所在的路径
	path := comm.GetCurrentPath()
	// 加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(path + "/config.ini")

	app := infra.New(conf)
	app.Start()
}
