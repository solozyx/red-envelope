package gorpc

import (
	"github.com/solozyx/red-envelope/services"
)

type EnvelopeRpc struct{}

// Go内置的 json rpc 接口规范
// 1. 入参和出参都要作为方法的参数
// 2. 方法必须有2个参数 并且是可导出类型
// 3. 方法的第二个参数(返回值)必须是一个指针
// 4. 方法返回值要返回一个error类型
// 5. 方法必须是可导出的

func (e *EnvelopeRpc) SendOut(in services.RedEnvelopeSendingDTO, out *services.RedEnvelopeActivity) error {
	s := services.GetRedEnvelopeService()
	a, err := s.SendOut(in)
	a.CopyTo(out)
	return err
}

func (e *EnvelopeRpc) Receive(in services.RedEnvelopeReceiveDTO, out *services.RedEnvelopeItemDTO) error {
	s := services.GetRedEnvelopeService()
	item, err := s.Receive(in)
	item.CopyTo(out)
	return err
}
