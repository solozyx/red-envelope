package accounts

import (
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
	_ "github.com/solozyx/red-envelope/textx"
)

func TestAccountLogDao(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountLogDao{
			runner: runner,
		}
		l := &AccountLog{
			LogNo:           ksuid.New().Next().String(),
			TradeNo:         ksuid.New().Next().String(),
			Status:          1,
			AccountNo:       ksuid.New().Next().String(),
			UserId:          ksuid.New().Next().String(),
			Username:        "账户流水测试",
			Amount:          decimal.NewFromFloat(1),
			Balance:         decimal.NewFromFloat(100),
			ChangeFlag:      services.FlagAccountCreated,
			ChangeType:      services.AccountCreated,
			TargetAccountNo: ksuid.New().Next().String(),
			TargetUserId:    ksuid.New().Next().String(),
			TargetUsername:  "目标账户流水测试",
			Desc:            "流水测试",
		}

		Convey("账户流水表测试", t, func() {
			// 通过log_no查询
			Convey("通过 log_no 查询", func() {
				id, err := dao.Insert(l)
				So(err, ShouldBeNil)
				So(id, ShouldBeGreaterThan, 0)

				out := dao.GetOne(l.LogNo)
				So(out, ShouldNotBeNil)
				So(out.Balance.String(), ShouldEqual, l.Balance.String())
				So(out.Amount.String(), ShouldEqual, l.Amount.String())
				So(out.CreatedAt, ShouldNotBeNil)
			})

			// 通过trade_no查询
			Convey("通过 trade_no 查询", func() {
				out := dao.GetByTradeNo(l.TradeNo)
				So(out, ShouldNotBeNil)
				So(out.Balance.String(), ShouldEqual, l.Balance.String())
				So(out.Amount.String(), ShouldEqual, l.Amount.String())
				So(out.CreatedAt, ShouldNotBeNil)
			})
		})

		return nil
	})

	if err != nil {
		logrus.Error(err)
	}
}
