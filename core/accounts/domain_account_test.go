package accounts

import (
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/solozyx/red-envelope/services"
)

func TestAccountDomain_Create(t *testing.T) {
	dto := services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "账户创建测试",
		Balance:  decimal.NewFromFloat(0),
		Status:   1,
	}
	// TODO:NOTICE 业务领域domain对象有状态每次使用都要重新初始化
	// domain := &accountDomain{}
	domain := new(accountDomain)

	Convey("账户创建测试", t, func() {
		rdto, err := domain.Create(dto)
		So(err, ShouldBeNil)
		So(rdto, ShouldNotBeNil)
		So(rdto.Balance.String(), ShouldEqual, dto.Balance.String())
		So(rdto.UserId, ShouldEqual, dto.UserId)
		So(rdto.Username, ShouldEqual, dto.Username)
		So(rdto.Status, ShouldEqual, dto.Status)
	})
}

func TestAccountDomain_Transfer(t *testing.T) {
	// 两个账户 一个交易主体 一个交易对象
	body := &services.AccountDTO{
		UserId:      ksuid.New().Next().String(),
		Username:    "交易主体",
		Balance:     decimal.NewFromFloat(100),
		Status:      1,
		AccountType: int(services.EnvelopeAccountType),
	}
	target := &services.AccountDTO{
		UserId:      ksuid.New().Next().String(),
		Username:    "交易对象",
		Balance:     decimal.NewFromFloat(100),
		Status:      1,
		AccountType: int(services.EnvelopeAccountType),
	}

	// 账户业务领域
	domain := accountDomain{}

	Convey("转账测试", t, func() {
		// 交易主体创建
		aBody, err := domain.Create(*body)
		So(err, ShouldBeNil)
		So(aBody, ShouldNotBeNil)
		So(aBody.Balance.String(), ShouldEqual, body.Balance.String())
		So(aBody.UserId, ShouldEqual, body.UserId)
		So(aBody.Username, ShouldEqual, body.Username)
		So(aBody.Status, ShouldEqual, body.Status)
		So(aBody.AccountName, ShouldEqual, body.AccountName)

		// 交易目标创建
		aTarget, err := domain.Create(*target)
		So(err, ShouldBeNil)
		So(aTarget, ShouldNotBeNil)
		So(aTarget.Balance.String(), ShouldEqual, target.Balance.String())
		So(aTarget.UserId, ShouldEqual, target.UserId)
		So(aTarget.Username, ShouldEqual, target.Username)
		So(aTarget.Status, ShouldEqual, target.Status)
		So(aTarget.AccountName, ShouldEqual, target.AccountName)

		// 转账操作验证
		Convey("余额充足，转入其他账户", func() {
			amount := decimal.NewFromFloat(1)
			dto := services.AccountTransferDTO{
				TradeNo: ksuid.New().Next().String(),
				TradeBody: services.TradeParticipator{
					AccountNo: aBody.AccountNo,
					UserId:    aBody.UserId,
					Username:  aBody.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: aTarget.AccountNo,
					UserId:    aTarget.UserId,
					Username:  aTarget.Username,
				},
				Amount:     amount,
				ChangeType: services.EnvelopeOutgoing,
				ChangeFlag: services.FlagTransferOut,
				Desc:       "余额充足，转入其他账户",
			}

			// 执行转账
			status, err := domain.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferredStatusSuccess)

			// 实际余额更新后的预期值验证
			account := domain.GetAccount(aBody.AccountNo)
			So(account, ShouldNotBeNil)
			So(account.Balance.String(), ShouldEqual, aBody.Balance.Sub(amount).String())
			So(err, ShouldBeNil)
		})

		Convey("余额不足，金额转出", func() {
			amount := aBody.Balance
			amount = amount.Add(decimal.NewFromFloat(200))
			dto := services.AccountTransferDTO{
				TradeNo: ksuid.New().Next().String(),
				TradeBody: services.TradeParticipator{
					AccountNo: aBody.AccountNo,
					UserId:    aBody.UserId,
					Username:  aBody.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: aTarget.AccountNo,
					UserId:    aTarget.UserId,
					Username:  aTarget.Username,
				},
				Amount:     amount,
				ChangeType: services.EnvelopeOutgoing,
				ChangeFlag: services.FlagTransferOut,
				Desc:       "余额不足，转入其他账户",
			}
			status, err := domain.Transfer(dto)
			// 扣减余额不足 应该返回错误
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, services.TransferredStatusSufficientFunds)

			// 实际余额更新后的预期值验证
			account := domain.GetAccount(aBody.AccountNo)
			So(account, ShouldNotBeNil)
			// 扣减失败 余额 等于 上次金额足够扣减后的金额
			So(account.Balance.String(), ShouldEqual, aBody.Balance.String())
			So(err, ShouldBeNil)
		})

		Convey("充值", func() {
			amount := decimal.NewFromFloat(11.1)
			dto := services.AccountTransferDTO{
				TradeNo: ksuid.New().Next().String(),
				TradeBody: services.TradeParticipator{
					AccountNo: aBody.AccountNo,
					UserId:    aBody.UserId,
					Username:  aBody.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: aBody.AccountNo,
					UserId:    aBody.UserId,
					Username:  aBody.Username,
				},
				Amount:     amount,
				ChangeType: services.AccountStoreValue,
				ChangeFlag: services.FlagTransferIn,
				Desc:       "充值",
			}
			status, err := domain.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferredStatusSuccess)

			// 实际余额更新后的预期值验证
			account := domain.GetAccount(aBody.AccountNo)
			So(account, ShouldNotBeNil)
			So(account.Balance.String(), ShouldEqual, aBody.Balance.Add(amount).String())
			So(err, ShouldBeNil)
		})
	})
}
