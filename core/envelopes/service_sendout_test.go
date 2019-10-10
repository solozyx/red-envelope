package envelopes

import (
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/solozyx/red-envelope/services"
	_ "github.com/solozyx/red-envelope/textx"
)

func TestRedEnvelopeService_SendOut(t *testing.T) {
	as := services.GetAccountService()
	accountDTO := services.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		AccountName:  "测试账户",
		AccountType:  int(services.EnvelopeAccountType),
		Amount:       "100",
		CurrencyCode: "CNY",
	}
	rs := services.GetRedEnvelopeService()

	Convey("准备账户", t, func() {
		account, err := as.CreateAccount(accountDTO)
		So(err, ShouldBeNil)
		So(account, ShouldNotBeNil)
		So(account.Balance.String(), ShouldEqual, accountDTO.Amount)
	})

	Convey("发红包测试", t, func() {
		goodsDTO := services.RedEnvelopeSendingDTO{
			Username: accountDTO.Username,
			UserId:   accountDTO.UserId,
			Blessing: "",
		}

		Convey("发普通红包", func() {
			goodsDTO.EnvelopeType = int(services.GeneralEnvelopeType)
			// 普通红包 用户输入单个红包金额 和 红包数量 数据库返回总金额
			goodsDTO.Amount = decimal.NewFromFloat(8.88)
			goodsDTO.Quantity = 10
			activity, err := rs.SendOut(goodsDTO)

			So(err, ShouldBeNil)
			So(activity, ShouldNotBeNil)
			So(activity.Link, ShouldNotBeEmpty)
			So(activity.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			// 验证每一个属性
			result := activity.RedEnvelopeGoodsDTO
			So(result.Username, ShouldEqual, goodsDTO.Username)
			So(result.UserId, ShouldEqual, goodsDTO.UserId)
			So(result.Amount.String(), ShouldEqual,
				goodsDTO.Amount.Mul(decimal.NewFromFloat(float64(goodsDTO.Quantity))).String())
			So(result.AmountOne.String(), ShouldEqual, goodsDTO.Amount.String())
		})

		Convey("发碰运气红包", func() {
			goodsDTO.EnvelopeType = int(services.LuckyEnvelopeType)
			// 碰运气红包 用户输入红包总金额 和 红包数量 数据库返回总金额
			goodsDTO.Amount = decimal.NewFromFloat(88.8)
			goodsDTO.Quantity = 10
			activity, err := rs.SendOut(goodsDTO)

			So(activity, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(activity.Link, ShouldNotBeEmpty)
			So(activity.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			// 验证每一个属性
			result := activity.RedEnvelopeGoodsDTO
			So(result.Username, ShouldEqual, goodsDTO.Username)
			So(result.UserId, ShouldEqual, goodsDTO.UserId)
			So(result.Amount.String(), ShouldEqual, goodsDTO.Amount.String())
			So(result.AmountOne.String(), ShouldEqual, decimal.NewFromFloat(0).String())
		})
	})
}
