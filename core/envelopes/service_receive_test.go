package envelopes

import (
	"strconv"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/solozyx/red-envelope/services"
)

func TestGoodsDomain_Receive(t *testing.T) {
	// 1.准备几个红包资金账户收发红包
	as := services.GetAccountService()
	accounts := make([]*services.AccountDTO, 0)
	size := 10

	Convey("收红包测试用例", t, func() {
		// 准备10个红包账户
		for i := 0; i < size; i++ {
			account := services.AccountCreatedDTO{
				UserId:       ksuid.New().Next().String(),
				Username:     "测试用户" + strconv.Itoa(i+1),
				AccountName:  "测试账户" + strconv.Itoa(i+1),
				AccountType:  int(services.EnvelopeAccountType),
				Amount:       "2000",
				CurrencyCode: "CNY",
			}
			// 账户创建
			acDto, err := as.CreateAccount(account)
			So(err, ShouldBeNil)
			So(acDto, ShouldNotBeNil)
			accounts = append(accounts, acDto)
		}

		// 2.选择第1个账户发红包
		acDTO := accounts[0]
		rs := services.GetRedEnvelopeService()
		// 发送普通红包
		goods := services.RedEnvelopeSendingDTO{
			UserId:       acDTO.UserId,
			Username:     acDTO.Username,
			EnvelopeType: int(services.GeneralEnvelopeType),
			// 普通红包 Amount 是每个子红包金额 总金额 = 1.88 * 10 = 18.8 元
			Amount:   decimal.NewFromFloat(1.88),
			Quantity: size,
			Blessing: "发红包",
		}
		// 发红包返回红包活动
		activity, err := rs.SendOut(goods)
		So(err, ShouldBeNil)
		So(activity, ShouldNotBeNil)
		So(activity.Link, ShouldNotBeEmpty)
		So(activity.RedEnvelopeGoodsDTO, ShouldNotBeNil)
		// 验证每一个属性
		dto := activity.RedEnvelopeGoodsDTO
		So(dto.UserId, ShouldEqual, goods.UserId)
		So(dto.Username, ShouldEqual, goods.Username)
		So(dto.Quantity, ShouldEqual, goods.Quantity)
		So(dto.Amount.String(), ShouldEqual, goods.Amount.Mul(decimal.NewFromFloat(float64(goods.Quantity))).String())
		So(dto.AmountOne.String(), ShouldEqual, goods.Amount.String())

		// 发红包后 剩余金额 = 总金额
		remainAmount := activity.Amount

		// 3.使用发送红包数量的人收红包 发红包的人也可以收红包
		Convey("收普通红包", func() {
			for _, account := range accounts {
				receiveDTO := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   activity.EnvelopeNo,
					RecvUserId:   account.UserId,
					RecvUsername: account.Username,
					AccountNo:    account.AccountNo,
				}
				// 收红包 返回红包订单详情
				item, err := rs.Receive(receiveDTO)
				So(err, ShouldBeNil)
				So(item, ShouldNotBeNil)
				// 收到的红包金额
				So(item.Amount.String(), ShouldEqual, activity.AmountOne.String())
				// 每次收红包后 红包剩余金额
				remainAmount = remainAmount.Sub(item.Amount)
				So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())
			}
		})

		Convey("收碰运气红包", func() {
			goods.EnvelopeType = int(services.LuckyEnvelopeType)
			at, err := rs.SendOut(goods)
			So(at, ShouldNotBeNil)
			So(err, ShouldBeNil)
			remainAmount = at.RemainAmount
			accounts = accounts[:10]
			for _, account := range accounts {
				receiveDTO := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   at.EnvelopeNo,
					RecvUsername: account.Username,
					RecvUserId:   account.UserId,
					AccountNo:    account.AccountNo,
				}
				item, err := rs.Receive(receiveDTO)
				So(err, ShouldBeNil)
				So(item, ShouldNotBeNil)
				remainAmount = remainAmount.Sub(item.Amount)
				So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())
			}
		})

	})
}
