package envelopes

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
)

const (
	// 分页 每页大小
	pageSize = 100
)

type ExpiredEnvelopeDomain struct {
	// 过期待退款红包列表
	expiredGoods []RedEnvelopeGoods
	// 分页 偏移量 从 offset 行开始 查询 pageSize 个数据
	offset int
}

// 查询出过期红包
func (e *ExpiredEnvelopeDomain) Next() (ok bool) {
	base.Tx(func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeGoodsDao{runner: runner}
		e.expiredGoods = dao.FindExpired(e.offset, pageSize)
		logrus.Infof("查询到 %d 个可退款红包", len(e.expiredGoods))
		if len(e.expiredGoods) > 0 {
			e.offset += len(e.expiredGoods)
			ok = true
		}
		return nil
	})
	return ok
}

func (e *ExpiredEnvelopeDomain) Expired() (err error) {
	for e.Next() {
		for _, g := range e.expiredGoods {
			logrus.Debugf("过期红包退款开始: %+v", g)
			err = e.ExpiredOne(g)
			if err != nil {
				logrus.Error(err)
			}
			logrus.Debugf("过期红包退款结束: %+v", g)
		}
	}
	return err
}

// 针对1个红包 发起1个退款流程
func (e *ExpiredEnvelopeDomain) ExpiredOne(goods RedEnvelopeGoods) (err error) {
	// 创建1个退款订单
	refund := goods
	refund.OrderType = services.OrderTypeRefund
	refund.Status = services.OrderExpired
	refund.PayStatus = services.Refunding
	refund.OriginEnvelopeNo = goods.EnvelopeNo
	refund.EnvelopeNo = ""
	// 过期时间 默认24hour
	refund.ExpiredAt = time.Now().Add(24 * time.Hour)
	domain := goodsDomain{RedEnvelopeGoods: refund}
	// 退款订单 红包商品生成新的红包编号 和 原过期红包编号 区分开
	domain.createEnvelopeNo()

	err = base.Tx(func(runner *dbx.TxRunner) error {
		txCtx := base.WithValueContext(context.Background(), runner)
		// 退款订单红包商品 写入 red_envelope_goods 表
		id, err := domain.Save(txCtx)
		if err != nil || id <= 0 {
			return errors.New("创建退款订单失败")
		}

		// 修改原过期红包订单状态
		dao := RedEnvelopeGoodsDao{runner: runner}
		_, err = dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpired)
		if err != nil {
			return errors.New("更新原过期红包订单为过期状态失败" + err.Error())
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 调用资金账户接口退款转账 系统红包账户 --> 原过期红包发送者账户
	systemAccount := base.GetSystemAccount()
	account := services.GetAccountService().GetEnvelopeAccountByUserId(goods.UserId)
	if account == nil {
		return errors.New("没有找到该用户的红包资金账户:" + goods.UserId)
	}
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: account.AccountNo,
		UserId:    account.UserId,
		Username:  account.Username,
	}

	// 系统账户扣减资金 转入原发红包账户
	transfer := services.AccountTransferDTO{
		TradeNo:     domain.RedEnvelopeGoods.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		Amount:      goods.RemainAmount, // 剩余金额转给红包发送人
		AmountStr:   goods.RemainAmount.String(),
		ChangeType:  services.SysEnvelopeExpiredRefund,
		ChangeFlag:  services.FlagTransferOut,
		Desc:        "过期红包退款,系统账户扣减资金,转给原红包发送人账户,红包编号: " + goods.EnvelopeNo,
	}
	status, err := services.GetAccountService().Transfer(transfer)
	if status != services.TransferredStatusSuccess {
		return err
	}

	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		// 修改原过期红包订单状态
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpiredRefundSucceed)
		if err != nil || rows == 0 {
			return errors.New("更新原过期红包订单状态为退款成功状态失败")
		}
		// 修改退款订单状态
		_, err = dao.UpdateOrderStatus(refund.EnvelopeNo, services.OrderExpiredRefundSucceed)
		if err != nil || rows == 0 {
			return errors.New("更退款订单状态为退款成功状态失败")
		}
		return nil
	})

	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
