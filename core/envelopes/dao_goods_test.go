package envelopes

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
	_ "github.com/solozyx/red-envelope/textx"
)

// 查询
func TestRedEnvelopeDao_GetOne(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		Convey("通过编号查询红包数据", t, func() {
			var good = &RedEnvelopeGoods{
				EnvelopeNo:     ksuid.New().Next().String(),
				EnvelopeType:   int(services.LuckyEnvelopeType),
				Username:       sql.NullString{String: ksuid.New().Next().String(), Valid: true},
				UserId:         ksuid.New().Next().String(),
				Blessing:       sql.NullString{String: "测试用红包商品", Valid: true},
				Amount:         decimal.NewFromFloat(100),
				AmountOne:      decimal.Decimal{},
				Quantity:       10,
				RemainAmount:   decimal.NewFromFloat(100),
				RemainQuantity: 10,
				ExpiredAt:      time.Now().Add(5 * time.Second),
				Status:         services.OrderSending,
				OrderType:      services.OrderTypeSending,
				PayStatus:      services.Payed,
			}
			id, err := dao.Insert(good)
			So(id, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)

			out := dao.GetOne(good.EnvelopeNo)
			So(out, ShouldNotBeNil)
			So(out.EnvelopeNo, ShouldEqual, good.EnvelopeNo)
			So(out.Username.String, ShouldEqual, good.Username.String)
			So(out.UserId, ShouldEqual, good.UserId)
			So(out.Amount.String(), ShouldEqual, good.Amount.String())
			So(out.Quantity, ShouldEqual, good.Quantity)
			So(out.RemainAmount.String(), ShouldEqual, good.RemainAmount.String())
			So(out.RemainQuantity, ShouldEqual, good.RemainQuantity)
			So(out.Status, ShouldEqual, good.Status)
			So(out.OrderType, ShouldEqual, good.OrderType)
			So(out.PayStatus, ShouldEqual, good.PayStatus)

		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

// 更新余额
func TestRedEnvelopeDao_UpdateBalance(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		Convey("通过编号查询红包数据", t, func() {
			var good = &RedEnvelopeGoods{
				EnvelopeNo:     ksuid.New().Next().String(),
				EnvelopeType:   int(services.LuckyEnvelopeType),
				Username:       sql.NullString{String: ksuid.New().Next().String(), Valid: true},
				UserId:         ksuid.New().Next().String(),
				Blessing:       sql.NullString{String: "测试用红包商品", Valid: true},
				Amount:         decimal.NewFromFloat(100),
				AmountOne:      decimal.Decimal{},
				Quantity:       10,
				RemainAmount:   decimal.NewFromFloat(100),
				RemainQuantity: 10,
				ExpiredAt:      time.Now().Add(5 * time.Second),
				Status:         services.OrderSending,
				OrderType:      services.OrderTypeSending,
				PayStatus:      services.Payed,
			}
			id, err := dao.Insert(good)
			So(id, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)

			amount := decimal.NewFromFloat(12)
			rows, err := dao.UpdateBalance(good.EnvelopeNo, amount)
			So(rows, ShouldEqual, 1)
			So(err, ShouldBeNil)

			out := dao.GetOne(good.EnvelopeNo)
			So(out, ShouldNotBeNil)
			So(out.EnvelopeNo, ShouldEqual, good.EnvelopeNo)
			So(out.Username.String, ShouldEqual, good.Username.String)
			So(out.UserId, ShouldEqual, good.UserId)
			So(out.Amount.String(), ShouldEqual, good.Amount.String())
			So(out.Quantity, ShouldEqual, good.Quantity)
			So(out.RemainAmount.String(), ShouldEqual, good.RemainAmount.Sub(amount).String())
			So(out.RemainQuantity, ShouldEqual, good.RemainQuantity-1)
			So(out.Status, ShouldEqual, good.Status)
			So(out.OrderType, ShouldEqual, good.OrderType)
			So(out.PayStatus, ShouldEqual, good.PayStatus)
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

// 过期查询
func TestRedEnvelopeDao_FindExpired(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		Convey("寻找过期", t, func() {
			goods := dao.FindExpired(0, 10)
			So(len(goods), ShouldBeGreaterThan, 0)
			for _, good := range goods {
				fmt.Println(good.EnvelopeNo)
			}

		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

// 更新状态
func TestRedEnvelopeDao_UpdateOrderStatus(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		Convey("更新状态", t, func() {
			// 插入数据
			var good = &RedEnvelopeGoods{
				EnvelopeNo:     ksuid.New().Next().String(),
				EnvelopeType:   int(services.LuckyEnvelopeType),
				Username:       sql.NullString{String: ksuid.New().Next().String(), Valid: true},
				UserId:         ksuid.New().Next().String(),
				Blessing:       sql.NullString{String: "测试用红包商品", Valid: true},
				Amount:         decimal.NewFromFloat(100),
				AmountOne:      decimal.Decimal{},
				Quantity:       10,
				RemainAmount:   decimal.NewFromFloat(100),
				RemainQuantity: 10,
				ExpiredAt:      time.Now().Add(5 * time.Second),
				Status:         services.OrderSending,
				OrderType:      services.OrderTypeSending,
				PayStatus:      services.Payed,
			}
			id, err := dao.Insert(good)
			So(id, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)

			rows, err := dao.UpdateOrderStatus(good.EnvelopeNo, 100)
			So(rows, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
