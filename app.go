package app

import (
	"github.com/solozyx/red-envelope/apis/gorpc"
	_ "github.com/solozyx/red-envelope/apis/gorpc"
	_ "github.com/solozyx/red-envelope/apis/web"
	_ "github.com/solozyx/red-envelope/core/accounts"
	_ "github.com/solozyx/red-envelope/core/envelopes"
	"github.com/solozyx/red-envelope/infra"
	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/jobs"
	_ "github.com/solozyx/red-envelope/views"
)

func init() {
	// 注册 配置文件读取启动器
	infra.Register(&base.PropsStarter{})
	// 注册 数据库启动
	infra.Register(&base.DbxDatabaseStarter{})
	// 注册 用户请求参数验证启动器
	infra.Register(&base.ValidatorStarter{})
	// 注册 RPC server
	infra.Register(&base.GoRPCStarter{})
	infra.Register(&gorpc.GoRPCApiStarter{})
	// 注册 过期红包退款 定时任务 要放在数据库starter之后 web starter 之前
	infra.Register(&jobs.RefundExpiredJobStarter{})
	infra.Register(&base.HookStarter{})

	// 注册 iris web server 是阻塞式放到最后位置
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
}
