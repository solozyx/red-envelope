package main

import (
	"fmt"

	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	// 和 locales/zh 包名冲突要用包别名
	vtzh "gopkg.in/go-playground/validator.v9/translations/zh"
)

type Person struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Age       uint8  `validate:"gte=0,lte=130"`
	Email     string `validate:"required,email"`
}

func main() {
	person := &Person{
		FirstName: "firstName",
		LastName:  "lastName",
		Age:       136,
		Email:     "fl163.com",
	}
	validate := validator.New()
	// 创建消息国际化通用翻译器
	// 中文翻译器
	cn := zh.New()
	// 通用翻译器
	uni := ut.New(cn, cn)
	// 创建本地化翻译器
	translator, found := uni.GetTranslator("zh")
	if found {
		// 验证器 注册 翻译器
		err := vtzh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("not found")
	}

	err := validate.Struct(person)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			fmt.Println(err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				// 使用注册的翻译器翻译错误信息
				fmt.Println(err.Translate(translator))
			}

		}
	}
}
