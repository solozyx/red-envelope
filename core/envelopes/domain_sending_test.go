package envelopes

import (
	"path"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/core/accounts"
	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
)

func TestGoodsDomain_SendOut(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		domain := &goodsDomain{}
		adomain := accounts.NewAccountDomain()
		aDto := &services.AccountDTO{
			AccountName: "测试用账户",
			AccountType: 1,
			UserId:      "发送测试用userId",
			Username:    "发送测试用username",
			Balance:     decimal.NewFromFloat(10000),
		}

		Convey("发送红包", t, func() {
			// 创建账户
			rdto, err := adomain.Create(aDto)
			So(rdto, ShouldNotBeNil)
			So(err, ShouldBeNil)
			dto := &services.RedEnvelopeGoodsDTO{
				EnvelopeType: services.LuckyEnvelopeType,
				Username:     "发送测试用username",
				UserId:       "发送测试用userId",
				Blessing:     "发送测试用祝福语",
				Amount:       decimal.NewFromFloat(1000),
				Quantity:     10,
				OrderType:    services.OrderTypeSending,
				AccountNo:    adomain.GetAccountNo(),
			}
			// 从该账户转账
			a, err := domain.SendOut(dto)
			So(a, ShouldNotBeNil)
			So(err, ShouldBeNil)

			So(a.Link, ShouldEqual, path.Join(base.GetEnvelopeDomain(), base.GetEnvelopeActivityLink(), domain.EnvelopeNo))
			So(a.EnvelopeNo, ShouldEqual, domain.EnvelopeNo)
			So(a.Blessing, ShouldEqual, domain.Blessing.String)
			So(a.Username, ShouldEqual, domain.Username.String)
			So(a.UserId, ShouldEqual, domain.UserId)
			So(a.Amount.String(), ShouldEqual, domain.Amount.String())
			So(a.AmountOne.String(), ShouldEqual, decimal.NewFromFloat(0).String())
			So(a.RemainAmount.String(), ShouldEqual, domain.RemainAmount.String())
			So(a.RemainQuantity, ShouldEqual, domain.RemainQuantity)
			So(a.Quantity, ShouldEqual, domain.Quantity)
		})

		return nil
	})
	if err != nil {
		logrus.Panic(err)
	}
}
