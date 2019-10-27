package gorpc

import (
	"github.com/solozyx/red-envelope/infra"
	"github.com/solozyx/red-envelope/infra/base"
)

type GoRPCApiStarter struct {
	infra.BaseStarter
}

func (s *GoRPCApiStarter) Init(ctx infra.StarterContext) {
	base.RpcRegister(new(EnvelopeRpc))
}
