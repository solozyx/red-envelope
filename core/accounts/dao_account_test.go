package accounts

import (
	"database/sql"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	// 静态导入 goconvey 直接使用 convey 包
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/infra/base"
	// 匿名导入
	_ "github.com/solozyx/red-envelope/textx"
)

func TestAccountDao_GetOne(t *testing.T) {
	txFunc := func(runner *dbx.TxRunner) error {
		dao := &AccountDao{
			runner: runner,
		}
		testFunc := func() {
			account := &Account{
				// 随机生成1个uid
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserId:      ksuid.New().Next().String(),
				Username:    sql.NullString{String: "测试资金用户", Valid: true},
				Balance:     decimal.NewFromFloat(1000),
				Status:      1,
			}
			// 账户创建
			id, err := dao.Insert(account)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
			// 账户查询
			out := dao.GetOne(account.AccountNo)
			So(out, ShouldNotBeNil)
			// goconvey不支持decimal类型 转换为字符串进行判断
			// decimal 认为 0.10  0.1000 是不同的值 转换为string会格式化 忽略小数点末尾的0
			So(out.Balance.String(), ShouldEqual, account.Balance.String())
			So(out.CreatedAt, ShouldNotBeNil)
			So(out.UpdatedAt, ShouldNotBeNil)
		}
		Convey("通过 创建账户 账户编号查询账户数据 ", t, testFunc)
		return nil
	}

	// 开始数据库操作事务
	err := base.Tx(txFunc)
	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_GetByUserId(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{
			runner: runner,
		}
		Convey("通过用户id和账户类型查询账户数据", t, func() {
			account := &Account{
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserId:      ksuid.New().Next().String(),
				Username: sql.NullString{
					String: "测试资金用户",
					Valid:  true,
				},
				Balance:     decimal.NewFromFloat(1000),
				Status:      1,
				AccountType: 2,
			}
			id, err := dao.Insert(account)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 1)
			out := dao.GetByUserId(account.UserId, account.AccountType)
			So(out, ShouldNotBeNil)
			So(out.Balance.String(), ShouldEqual, account.Balance.String())
			So(out.CreatedAt, ShouldNotBeNil)
			So(out.UpdatedAt, ShouldNotBeNil)
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

// 账户余额扣减
func TestAccountDao_UpdateBalance(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{
			runner: runner,
		}
		balance := decimal.NewFromFloat(100)
		Convey("更新账户余额", t, func() {
			account := &Account{
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserId:      ksuid.New().Next().String(),
				Username: sql.NullString{
					String: "测试资金用户",
					Valid:  true,
				},
				Balance:     balance,
				Status:      1,
				AccountType: 1,
			}
			id, err := dao.Insert(account)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 1)

			// 1.增加余额
			Convey("增加余额", func() {
				amount := decimal.NewFromFloat(10)
				rows, err := dao.UpdateBalance(account.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 1)

				newBalance := balance.Add(amount)
				out := dao.GetOne(account.AccountNo)
				So(out, ShouldNotBeNil)
				So(newBalance.String(), ShouldEqual, out.Balance.String())
				So(out.CreatedAt, ShouldNotBeNil)
			})

			// 2.扣减余额 余额足够
			Convey("扣减余额，余额足够", func() {
				amount := decimal.NewFromFloat(-100)
				rows, err := dao.UpdateBalance(account.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 1)

				newBalance := balance.Add(amount)
				out := dao.GetOne(account.AccountNo)
				So(out, ShouldNotBeNil)
				So(newBalance.String(), ShouldEqual, out.Balance.String())
				So(out.CreatedAt, ShouldNotBeNil)
			})

			// 3.扣减余额 余额不够
			Convey("扣减余额，余额不足", func() {
				out1 := dao.GetOne(account.AccountNo)
				So(out1, ShouldNotBeNil)

				amount := decimal.NewFromFloat(-100000)
				rows, err := dao.UpdateBalance(account.AccountNo, amount)
				// 受影响行数 扣减失败 0行
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 0)

				out2 := dao.GetOne(account.AccountNo)
				So(out2, ShouldNotBeNil)
				So(out1.Balance.String(), ShouldEqual, out2.Balance.String())
				So(out1.CreatedAt, ShouldEqual, out2.CreatedAt)
			})
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
