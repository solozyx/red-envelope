package base

import "github.com/sirupsen/logrus"

func Check(i interface{}) {
	if i == nil {
		logrus.Panicf("参数未被实例化, %v", i)
		panic("参数未被实例化")
	}
}
