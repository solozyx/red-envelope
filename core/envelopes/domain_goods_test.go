package envelopes

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
)

func TestGoodsDomain_CreateAndSave(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		domain := &goodsDomain{}
		dao := &RedEnvelopeDao{runner: runner}
		dto := &services.RedEnvelopeGoodsDTO{
			EnvelopeType: services.LuckyEnvelopeType,
			Username:     "测试用username",
			UserId:       "测试用userId",
			Blessing:     "测试用祝福语",
			Amount:       decimal.NewFromFloat(1000),
			Quantity:     10,
			OrderType:    services.OrderTypeSending,
			AccountNo:    "测试用账户",
		}
		Convey("创建和保存红包商品", t, func() {
			ctx := base.WithValueContext(context.Background(), runner)
			id, err := domain.CreateAndSave(ctx, dto)
			So(id, ShouldBeGreaterThan, 1)
			So(err, ShouldBeNil)

			out := dao.GetOne(domain.EnvelopeNo)
			So(out, ShouldNotBeNil)
			So(out.EnvelopeNo, ShouldEqual, domain.EnvelopeNo)
			So(out.Amount.String(), ShouldEqual, domain.Amount.String())
			So(out.AmountOne.String(), ShouldEqual, decimal.NewFromFloat(0).String())
			So(out.EnvelopeType, ShouldEqual, domain.EnvelopeType)
			So(out.Blessing.String, ShouldEqual, domain.Blessing.String)
			So(out.Blessing.Valid, ShouldEqual, domain.Blessing.Valid)
			So(out.Quantity, ShouldEqual, domain.Quantity)
			So(out.RemainQuantity, ShouldEqual, domain.RemainQuantity)
			So(out.RemainAmount.String(), ShouldEqual, domain.RemainAmount.String())
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
