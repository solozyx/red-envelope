package textx

import (
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/tietang/props/ini"

	"github.com/solozyx/red-envelope/comm"
	"github.com/solozyx/red-envelope/infra"
	"github.com/solozyx/red-envelope/infra/base"
)

func init() {
	// 获取程序运行文件所在的路径
	path := comm.GetCurrentPath()
	path = strings.TrimRight(path, "textx")
	path += "brun/config.ini"
	logrus.Info("配置文件路径 %s", path)
	// 加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(path)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})

	app := infra.New(conf)
	app.Start()
}
