package infra

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
)

const (
	KeyProps = "_conf"
)

// 基础资源上下文
type StarterContext map[string]interface{}

func (s StarterContext) Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置还没有被初始化")
	}
	return p.(kvs.ConfigSource)
}

// 资源启动器接口
type Starter interface {
	// 1.系统启动，初始化一些基础资源
	Init(StarterContext)
	// 2.系统基础资源的安装
	Setup(StarterContext)
	// 3.启动基础资源
	Start(StarterContext)
	// 启动器是否可阻塞
	StartBlocking() bool
	// 4.资源停止和销毁
	Stop(StarterContext)
}

// 简单的验证 BaseStarter 是否实现了Starter接口
var _ Starter = new(BaseStarter)

// 基础空启动器实现 为了方便资源启动器的代码实现
type BaseStarter struct{}

func (b *BaseStarter) Init(ctx StarterContext)  {}
func (b *BaseStarter) Setup(ctx StarterContext) {}
func (b *BaseStarter) Start(ctx StarterContext) {}
func (b *BaseStarter) StartBlocking() bool      { return false }
func (b *BaseStarter) Stop(ctx StarterContext)  {}

// 启动器注册器
type starterRegister struct {
	nonBlockingStarters []Starter
	blockingStarters    []Starter
}

// 注册启动器
func (r *starterRegister) Register(starter Starter) {
	if starter.StartBlocking() {
		r.blockingStarters = append(r.blockingStarters, starter)
	} else {
		r.nonBlockingStarters = append(r.nonBlockingStarters, starter)
	}
	typ := reflect.TypeOf(starter)
	logrus.Infof("Register starter: %s", typ.String())
}

// 返回所有 Starter
func (r *starterRegister) AllStarters() []Starter {
	starters := make([]Starter, 0)
	starters = append(starters, r.nonBlockingStarters...)
	starters = append(starters, r.blockingStarters...)
	return starters
}

var StarterRegister *starterRegister = new(starterRegister)

func GetStarters() []Starter {
	return StarterRegister.AllStarters()
}

func Register(s Starter) {
	StarterRegister.Register(s)
}

// 系统基础资源的启动管理
func SystemRun() {
	ctx := StarterContext{}
	// 1.初始化
	for _, s := range GetStarters() {
		s.Init(ctx)
	}
	// 2.安装
	for _, s := range GetStarters() {
		s.Setup(ctx)
	}
	// 3.启动
	for _, s := range GetStarters() {
		s.Start(ctx)
	}
}
