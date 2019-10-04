package base

import (
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	// 和内置函数 recover 冲突 使用包别名
	irisrecover "github.com/kataras/iris/middleware/recover"
	"github.com/sirupsen/logrus"

	"github.com/solozyx/red-envelope/infra"
)

var irisApplication *iris.Application

func Iris() *iris.Application {
	return irisApplication
}

type IrisServerStarter struct {
	infra.BaseStarter
}

func (i *IrisServerStarter) Init(ctx infra.StarterContext) {
	logrus.Info("IrisServerStarter Init()")
	// 创建 iris application 实例
	irisApplication = initIris()
	// 日志组件的配置和扩展 iris内部使用它自己扩展的 golog 框架和 logrus 输出格式不同
	// 程序日志后期分析,输出日志格式不同则增加额外难度
	// 统一 iris golog 和 logrus 日志输出格式
	// 获取 iris 日志组件
	irisLogger := irisApplication.Logger()
	// Install方法扩展 golog 日志组件与logrus 统一 golog和logrus都实现了这套接口 所以可以适配
	irisLogger.Install(logrus.StandardLogger())
}

func (i *IrisServerStarter) Start(ctx infra.StarterContext) {
	// 把路由信息打印到控制台方便调试
	routers := Iris().GetRoutes()
	for _, r := range routers {
		// 打印路由信息
		logrus.Info(r.Trace())
	}
	// 启动 iris web server
	// 配置文件读取启动端口
	port := ctx.Props().GetDefault("app.server.port", "18080")
	// 监听所有网卡
	Iris().Run(iris.Addr(":" + port))
}

// IrisServerStarter 是阻塞式的 需要实现阻塞接口
func (i *IrisServerStarter) StartBlocking() bool {
	return true
}

func initIris() *iris.Application {
	app := iris.New()
	// 主要中间件的配置 recover 日志输出中间件的自定义
	// iris 内置 recover 中间件
	app.Use(irisrecover.New())
	// iris 内置 日志中间件
	cfg := logger.Config{
		// 状态
		Status: true,
		// IP地址
		IP: true,
		// http请求方式
		Method: true,
		// http请求path
		Path: true,
		// http请求query参数
		Query:              true,
		Columns:            false,
		MessageContextKeys: nil,
		MessageHeaderKeys:  nil,
		// 日志格式化输出函数
		LogFunc: func(
			// 时间戳
			now time.Time,
			// 延迟 请求响应时间
			latency time.Duration,
			// 状态 ip地址 http访问方式 http请求path
			status, ip, method, path string,
			// 消息
			message interface{},
			// 头消息
			headerMessage interface{}) {
			app.Logger().Infof("| %s | %s | %s | %s | %s | %s | %s | %s |",
				// 服务请求时间戳格式 和 日志本身打印时间 有区别
				now.Format("2006-01-02 15:04:05.000000"),
				// 请求延迟
				latency.String(),
				status, ip, method, path, headerMessage, message)
		},
		Skippers: nil,
	}
	// 添加日志中间件
	app.Use(logger.New(cfg))
	return app
}
