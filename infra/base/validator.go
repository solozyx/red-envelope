package base

import (
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	vtzh "gopkg.in/go-playground/validator.v9/translations/zh"

	"github.com/solozyx/red-envelope/infra"
)

var (
	// 全局验证器
	validate *validator.Validate
	// 全局翻译器
	translator ut.Translator
)

func Validate() *validator.Validate {
	return validate
}

func Translate() ut.Translator {
	return translator
}

type ValidatorStarter struct {
	infra.BaseStarter
}

// 验证器翻译器 和 其他基础组件依赖性不强
// 只需要在接收用户请求之前创建好即可,即只需要先于web服务器创建前注册好验证器即可
func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	logrus.Info("ValidatorStarter Init()")
	validate = validator.New()
	// 创建消息国际化通用翻译器
	cn := zh.New()
	uni := ut.New(cn, cn)
	var found bool
	translator, found = uni.GetTranslator("zh")
	if found {
		err := vtzh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		logrus.Error("Not found translator: zh")
	}
}

func ValidateStruct(s interface{}) error {
	err := Validate().Struct(s)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("账户创建验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				logrus.Error(e.Translate(Translate()))
			}
		}
		return err
	}
	return nil
}
