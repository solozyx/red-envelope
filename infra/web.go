package infra

var apiInitializerRegister = new(InitializerRegister)

// 注册webApi初始化对象
func RegisterApi(init Initializer) {
	apiInitializerRegister.register(init)
}

func GetApiInitializers() []Initializer {
	return apiInitializerRegister.Initializers
}

type WebApiStarter struct {
	BaseStarter
}

func (s *WebApiStarter) Setup(ctx StarterContext) {
	for _, i := range GetApiInitializers() {
		i.Init()
	}
}
