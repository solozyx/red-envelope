package accounts

import (
	"github.com/shopspring/decimal"
	"testing"

	"github.com/segmentio/ksuid"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/solozyx/red-envelope/services"
)

func TestAccountService_CreateAccount(t *testing.T) {
	dto := services.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "账户创建测试",
		AccountName:  "账户创建测试",
		AccountType:  int(services.EnvelopeAccountType),
		CurrencyCode: "CNY",
		Amount:       "100",
	}
	s := new(accountService)
	Convey("账户创建测试", t, func() {
		rdto, err := s.CreateAccount(dto)
		So(err, ShouldBeNil)
		So(rdto, ShouldNotBeNil)
		So(rdto.Balance.String(), ShouldEqual, dto.Amount)
		So(rdto.UserId, ShouldEqual, dto.UserId)
		So(rdto.Username, ShouldEqual, dto.Username)
		So(rdto.AccountName, ShouldEqual, dto.AccountName)
		So(rdto.CurrencyCode, ShouldEqual, dto.CurrencyCode)
		So(rdto.AccountType, ShouldEqual, dto.AccountType)
		So(rdto.Status, ShouldEqual, 1)
	})
}

// 转账服务应用层测试用例代码
func TestAccountService_Transfer(t *testing.T) {
	Convey("转账测试", t, func() {
		// 创建2个账户
		a1 := services.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "用户测试1",
			AccountName:  "账户测试1",
			AccountType:  int(services.EnvelopeAccountType),
			CurrencyCode: "CNY",
			Amount:       "100",
		}
		a2 := services.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "用户测试2",
			AccountName:  "账户测试2",
			AccountType:  int(services.EnvelopeAccountType),
			CurrencyCode: "CNY",
			Amount:       "100",
		}

		// 实例化 service
		s := new(accountService)

		// 创建账户1
		a1DTO, err := s.CreateAccount(a1)
		So(err, ShouldBeNil)
		So(a1DTO, ShouldNotBeNil)
		So(a1DTO.UserId, ShouldEqual, a1.UserId)
		So(a1DTO.Username, ShouldEqual, a1.Username)
		So(a1DTO.AccountName, ShouldEqual, a1.AccountName)
		So(a1DTO.AccountType, ShouldEqual, a1.AccountType)
		So(a1DTO.CurrencyCode, ShouldEqual, a1.CurrencyCode)
		So(a1DTO.Balance.String(), ShouldEqual, a1.Amount)
		So(a1DTO.Status, ShouldEqual, 1)

		// 创建账户2
		a2DTO, err := s.CreateAccount(a2)
		So(err, ShouldBeNil)
		So(a2DTO, ShouldNotBeNil)
		So(a2DTO.UserId, ShouldEqual, a2.UserId)
		So(a2DTO.Username, ShouldEqual, a2.Username)
		So(a2DTO.AccountName, ShouldEqual, a2.AccountName)
		So(a2DTO.AccountType, ShouldEqual, a2.AccountType)
		So(a2DTO.CurrencyCode, ShouldEqual, a2.CurrencyCode)
		So(a2DTO.Balance.String(), ShouldEqual, a2.Amount)
		So(a2DTO.Status, ShouldEqual, 1)

		Convey("从账户1转入账户2，其中账户1余额足够", func() {
			amount := decimal.NewFromFloat(10)
			tDTO := services.AccountTransferDTO{
				TradeNo: ksuid.New().Next().String(),
				TradeBody: services.TradeParticipator{
					AccountNo: a1DTO.AccountNo,
					UserId:    a1DTO.UserId,
					Username:  a1DTO.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: a2DTO.AccountNo,
					UserId:    a2DTO.UserId,
					Username:  a2DTO.Username,
				},
				// 转账金额用字符串传输 避免在计算过程中丢失精度
				AmountStr:  amount.String(),
				ChangeType: services.EnvelopeOutgoing,
				ChangeFlag: services.FlagTransferOut,
				Desc:       "从账户1转入账户2，其中账户1余额足够",
			}
			status, err := s.Transfer(tDTO)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferredStatusSuccess)

			// 验证账户1金额变化
			a1DTOByAccountNo := s.GetAccount(a1DTO.AccountNo)
			So(a1DTOByAccountNo, ShouldNotBeNil)
			So(a1DTOByAccountNo.Balance.String(), ShouldEqual, a1DTO.Balance.Sub(amount).String())

			// 验证账户2金额变化
			a2DTOByAccountNo := s.GetAccount(a2DTO.AccountNo)
			So(a2DTOByAccountNo, ShouldNotBeNil)
			So(a2DTOByAccountNo.Balance.String(), ShouldEqual, a2DTO.Balance.Add(amount).String())
		})

		Convey("从账户2转入账户1，其中账户2余额不足，转账应该失败", func() {
			amount := decimal.NewFromFloat(200)
			tDTO := services.AccountTransferDTO{
				TradeNo: ksuid.New().Next().String(),
				TradeBody: services.TradeParticipator{
					AccountNo: a2DTO.AccountNo,
					UserId:    a2DTO.UserId,
					Username:  a2DTO.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: a1DTO.AccountNo,
					UserId:    a1DTO.UserId,
					Username:  a1DTO.Username,
				},
				AmountStr:  amount.String(),
				ChangeType: services.EnvelopeOutgoing,
				ChangeFlag: services.FlagTransferOut,
				Desc:       "从账户2转入账户1，其中账户2余额不足",
			}
			status, err := s.Transfer(tDTO)
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, services.TransferredStatusSufficientFunds)

			// 验证账户1金额变化
			a1DTOByAccountNo := s.GetAccount(a1DTO.AccountNo)
			So(a1DTOByAccountNo, ShouldNotBeNil)
			So(a1DTOByAccountNo.Balance.String(), ShouldEqual, "90")

			// 验证账户2金额变化
			a2DTOByAccountNo := s.GetAccount(a2DTO.AccountNo)
			So(a2DTOByAccountNo, ShouldNotBeNil)
			So(a2DTOByAccountNo.Balance.String(), ShouldEqual, "110")
		})

		Convey("给账户1储值", func() {
			tDTO := services.AccountTransferDTO{
				TradeNo: ksuid.New().Next().String(),
				TradeBody: services.TradeParticipator{
					AccountNo: a1DTO.AccountNo,
					UserId:    a1DTO.UserId,
					Username:  a1DTO.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: a1DTO.AccountNo,
					UserId:    a1DTO.UserId,
					Username:  a1DTO.Username,
				},
				AmountStr:  "100",
				Desc:       "给账户1储值",
				ChangeType: services.AccountStoreValue,
				ChangeFlag: services.FlagTransferIn,
			}
			status, err := s.StoreValue(tDTO)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferredStatusSuccess)

			// 验证账户金额变化
			a1DTOByAccountNo := s.GetAccount(a1DTO.AccountNo)
			So(a1DTOByAccountNo, ShouldNotBeNil)
			So(a1DTOByAccountNo.Balance.String(), ShouldEqual, "190")
		})
	})
}
