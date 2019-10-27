package envelopes

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
)

type goodsDomain struct {
	RedEnvelopeGoods
	itemDomain
}

// 生成1个红包编号
func (domain *goodsDomain) createEnvelopeNo() {
	domain.RedEnvelopeGoods.EnvelopeNo = ksuid.New().Next().String()
}

// 创建1个红包商品对象
func (domain *goodsDomain) Create(dto services.RedEnvelopeGoodsDTO) {
	domain.RedEnvelopeGoods.FromDTO(&dto)
	// sql.NullString Valid设为true才能持久化到数据库
	domain.Username.Valid = true
	domain.Blessing.Valid = true
	// 红包商品创建 remain_amount = amount remain_quantity = quantity
	domain.RemainQuantity = domain.Quantity
	if domain.EnvelopeType == int(services.GeneralEnvelopeType) {
		// 普通红包总金额 = 用户输入金额是单个红包金额amountOne * quantity
		// domain.Amount = dto.AmountOne.Mul(decimal.NewFromFloat(float64(dto.Quantity)))
		amountOne, _ := decimal.NewFromString(dto.AmountOne)
		domain.RedEnvelopeGoods.AmountOne = amountOne
		domain.RedEnvelopeGoods.Amount = amountOne.Mul(decimal.NewFromFloat(float64(dto.Quantity)))
	}
	if domain.EnvelopeType == int(services.LuckyEnvelopeType) {
		// 碰运气红包 用户输入的红包金额 = 红包总金额
		domain.RedEnvelopeGoods.AmountOne = decimal.NewFromFloat(0)
		domain.RedEnvelopeGoods.Amount, _ = decimal.NewFromString(dto.Amount)
	}
	domain.RemainAmount = domain.Amount
	// 过期时间 默认24hour
	domain.ExpiredAt = time.Now().Add(24 * time.Hour)
	domain.Status = services.OrderCreate
	domain.PayStatus = services.Paying
	domain.createEnvelopeNo()
}

// 保存到红包商品表
func (domain *goodsDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		id, err = dao.Insert(&domain.RedEnvelopeGoods)
		return err
	})
	// 0 表示保存到数据库失败
	return id, err
}

// 创建并保存红包商品
func (domain *goodsDomain) CreateAndSave(ctx context.Context, dto services.RedEnvelopeGoodsDTO) (id int64, err error) {
	domain.Create(dto)
	return domain.Save(ctx)
}

// 查询红包商品信息
func (domain *goodsDomain) Get(envelopeNo string) (goods *RedEnvelopeGoods) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeGoodsDao{runner: runner}
		goods = dao.GetOne(envelopeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return goods
}
