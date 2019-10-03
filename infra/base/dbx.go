package base

import (
	// MySQL驱动
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/tietang/props/kvs"

	"github.com/solozyx/red-envelope/infra"
)

// dbx 数据库实例
var database *dbx.Database

func DbxDatabase() *dbx.Database {
	return database
}

// dbx数据库starter 设置为全局
type DbxDatabaseStarter struct {
	// 继承 BaseStarter 免去实现一系列的 start方法
	infra.BaseStarter
}

// TODO:NOTICE 数据库连接启动生命周期要晚于配置文件加载 数据库资源初始化放到 Setup阶段
func (s *DbxDatabaseStarter) Setup(ctx infra.StarterContext) {
	logrus.Info("DbxDatabaseStarter Setup()")
	// 读取配置文件
	conf := ctx.Props()
	// 数据库配置
	settings := dbx.Settings{}
	// kvs.Unmarshal 把配置文件内容解析到结构体
	err := kvs.Unmarshal(conf, &settings, "mysql")
	if err != nil {
		// 数据库配置错误 是严重错误 禁止启动成功
		panic(err)
	}
	logrus.Infof("%+v", settings)
	// 记录连接字符串
	logrus.Info("mysql.conn url:", settings.ShortDataSourceName())
	// 实例化数据库连接对象
	dbConn, err := dbx.Open(settings)
	if err != nil {
		// 数据库连接异常 禁止启动应用
		logrus.Panic("dbx.Setup dbx.Open error:", err)
	}
	logrus.Info(dbConn.Ping())
	database = dbConn
}
