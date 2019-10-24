package infra

import (
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
)

// 应用程序启动管理器
type BootApplication struct {
	conf           kvs.ConfigSource
	starterContext StarterContext
}

func New(conf kvs.ConfigSource) *BootApplication {
	b := &BootApplication{
		conf:           conf,
		starterContext: StarterContext{},
	}
	b.starterContext[KeyProps] = conf
	return b
}

func (b *BootApplication) Start() {
	// 1.初始化所有starter
	b.init()
	// 2.安装所有starter
	b.setup()
	// 3.启动所有starter
	b.start()
}

func (b *BootApplication) init() {
	for _, starter := range GetStarters() {
		starter.Init(b.starterContext)
	}
}

func (b *BootApplication) setup() {
	for _, starter := range GetStarters() {
		starter.Setup(b.starterContext)
	}
}

func (b *BootApplication) start() {
	log.Info("Starting starters...")
	for i, starter := range GetStarters() {
		typ := reflect.TypeOf(starter)
		log.Debug("Starting: ", typ.String())
		if starter.StartBlocking() {
			// 如果是最后一个可阻塞的 直接启动并阻塞
			if i+1 == len(GetStarters()) {
				starter.Start(b.starterContext)
			} else {
				// 如果不是 使用goroutine来异步启动 防止阻塞后面的starter
				go starter.Start(b.starterContext)
			}
		} else {
			// 非阻塞直接运行
			starter.Start(b.starterContext)
		}
	}
}

func (b *BootApplication) Stop() {
	log.Info("Stoping starters...")
	for _, v := range GetStarters() {
		typ := reflect.TypeOf(v)
		log.Debug("Stoping: ", typ.String())
		v.Stop(b.starterContext)
	}
}
