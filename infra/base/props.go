package base

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"

	"github.com/solozyx/red-envelope/infra"
)

var props kvs.ConfigSource

// 获取全局配置实例
func Props() kvs.ConfigSource {
	Check(props)
	return props
}

type PropsStarter struct {
	// 继承 BaseStarter
	infra.BaseStarter
}

// 配置相关启动优先级最高,程序启动要首先读取配置文件
func (p *PropsStarter) Init(ctx infra.StarterContext) {
	logrus.Info("PropsStarter Init()")
	// props = ini.NewIniFileConfigSource("config.ini")
	props = ctx.Props()
	// 配置系统红包账户
	GetSystemAccount()
	fmt.Println("初始化配置。")
}

// 系统红包账户
type SystemAccount struct {
	AccountNo   string
	AccountName string
	UserId      string
	Username    string
}

var systemAccount *SystemAccount

// 保证 systemAccount 只初始化1次
var systemAccountOnce sync.Once

func GetSystemAccount() *SystemAccount {
	systemAccountOnce.Do(func() {
		// 不论 Do 方法调用多少次,该匿名函数只会执行1次
		systemAccount = new(SystemAccount)
		// 反序列配置文件 config.int 内容 读取配置文件的key前缀 system.account
		err := kvs.Unmarshal(Props(), systemAccount, "system.account")
		if err != nil {
			// 发红包必须有红包账户 否则业务无法进行 panic
			logrus.Panic(err)
		}
	})
	return systemAccount
}

func GetEnvelopeActivityLink() string {
	// 读取配置文件
	return Props().GetDefault("envelope.link", "/v1/envelope/link")
}

func GetEnvelopeDomain() string {
	return Props().GetDefault("envelope.domain", "http://localhost")
}
