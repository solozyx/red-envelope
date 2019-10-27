package web

import (
	"github.com/kataras/iris"

	"github.com/solozyx/red-envelope/infra"
	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
)

type RedEnvelopeApi struct {
	service services.RedEnvelopeService
}

func init() {
	infra.RegisterApi(&RedEnvelopeApi{})
}

func (api *RedEnvelopeApi) Init() {
	api.service = services.GetRedEnvelopeService()
	groupRouter := base.Iris().Party("/v1/envelope")
	groupRouter.Post("/sendout", api.sendOutHandler)
	groupRouter.Post("/receive", api.receiveHandler)
}

/*
{
	"envelope_type":0,
	"username":"",
	"user_id":"",
	"amount":"0",
	"quantity":0
}
*/
func (api *RedEnvelopeApi) sendOutHandler(ctx iris.Context) {
	dto := services.RedEnvelopeSendingDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsErr
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	// 发红包
	activity, err := api.service.SendOut(dto)
	if err != nil {
		r.Code = base.ResCodeInternalServerErr
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = activity
	ctx.JSON(&r)
}

func (api *RedEnvelopeApi) receiveHandler(ctx iris.Context) {
	dto := services.RedEnvelopeReceiveDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsErr
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	// 收红包
	item, err := api.service.Receive(dto)
	if err != nil {
		r.Code = base.ResCodeInternalServerErr
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = item
	ctx.JSON(&r)
}
