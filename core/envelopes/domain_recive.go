package envelopes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/core/accounts"
	"github.com/solozyx/red-envelope/infra/algo"
	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
)

// 人民币CNY 元 -> 分
var multiple = decimal.NewFromFloat(100.0)

// 收红包业务
func (domain *goodsDomain) Receive(ctx context.Context, dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	// 1.创建收红包的订单明细
	domain.preCreateItem(dto)
	// 2.查询出当前红包的剩余数量和剩余金额信息
	goods := domain.Get(dto.EnvelopeNo)
	// 3.校验剩余红包数量和剩余金额 如果没有剩余 直接返回无可用红包金额
	if goods.RemainQuantity <= 0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0)) <= 0 {
		return nil, errors.New("收红包没有足够的剩余金额")
	}
	// 4.使用红包算法计算红包金额
	nextAmount := domain.nextAmount(goods)

	err = base.Tx(func(runner *dbx.TxRunner) error {
		// 5.使用乐观锁更新语句 尝试更新剩余数量和剩余金额
		// - 更新成功 返回1 抢到红包
		// - 更新失败 返回0 无剩余红包金额或数量 抢红包失败
		dao := &RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateBalance(goods.EnvelopeNo, nextAmount)
		// 如果更新失败 row affected 返回0 表示无可用红包数量与金额
		if rows <= 0 || err != nil {
			return errors.New("收红包更新数据库失败 导致收红包失败")
		}
		// 如果更新成功 row affected 返回1 表示抢到红包
		// 6.保存订单明细数据
		domain.itemDomain.RedEnvelopeItem.Quantity = 1
		domain.itemDomain.RedEnvelopeItem.PayStatus = int(services.Paying)
		domain.itemDomain.RedEnvelopeItem.AccountNo = dto.AccountNo
		// 本次抢到红包后的红包剩余金额
		domain.itemDomain.RedEnvelopeItem.RemainAmount = goods.RemainAmount.Sub(nextAmount)
		// 本次抢到的红包金额
		domain.itemDomain.RedEnvelopeItem.Amount = nextAmount
		// 构造新的上下文 传入 *TxRunner 数据库事务对象
		txCtx := base.WithValueContext(ctx, runner)
		// 写入 red_envelope_item 表
		_, err = domain.itemDomain.Save(txCtx)
		if err != nil {
			return err
		}
		// 7.将抢到的红包金额从系统红包中间账户转入当前抢红包用户的资金账户
		status, err := domain.transfer(txCtx, dto)
		if status == services.TransferredStatusSuccess {
			return nil
		}
		return err
	})
	return domain.itemDomain.RedEnvelopeItem.ToDTO(), err
}

// 预创建收红包订单明细
func (domain *goodsDomain) preCreateItem(dto services.RedEnvelopeReceiveDTO) {
	// 收红包 和 账户 关联
	domain.itemDomain.RedEnvelopeItem.AccountNo = dto.AccountNo
	// // 收红包 和 发出的红包编号 关联
	domain.itemDomain.RedEnvelopeItem.EnvelopeNo = dto.EnvelopeNo
	domain.itemDomain.RedEnvelopeItem.RecvUsername = sql.NullString{
		String: dto.RecvUsername,
		Valid:  true,
	}
	domain.itemDomain.RedEnvelopeItem.RecvUserId = dto.RecvUserId

	envelopeGoods := domain.Get(dto.EnvelopeNo)
	var s string
	if envelopeGoods.EnvelopeType == int(services.GeneralEnvelopeType) {
		s = "普通"
	} else {
		s = "碰运气"
	}
	domain.itemDomain.RedEnvelopeItem.Desc = fmt.Sprintf("%s的%s红包",
		envelopeGoods.Username.String, s)

	domain.itemDomain.createItemNo()
}

// 计算红包金额
func (domain *goodsDomain) nextAmount(goods *RedEnvelopeGoods) (amount decimal.Decimal) {
	if goods.RemainQuantity == 1 {
		return goods.RemainAmount
	}
	if goods.EnvelopeType == int(services.GeneralEnvelopeType) {
		return goods.AmountOne
	}
	if goods.EnvelopeType == int(services.LuckyEnvelopeType) {
		// 剩余金额 元 -> 分 *100 取出int值
		centInt := goods.RemainAmount.Mul(multiple).IntPart()
		next := algo.DoubleAverage(int64(goods.RemainQuantity), centInt)
		// 分 -> 元 /100
		amount = decimal.NewFromFloat(float64(next)).Div(multiple)
	}
	return amount
}

func (domain *goodsDomain) transfer(ctx context.Context, dto services.RedEnvelopeReceiveDTO) (status services.TransferredStatus, err error) {
	// 交易主体 系统红包账户
	systemAccount := base.GetSystemAccount()
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	// 交易对方
	target := services.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		Username:  dto.RecvUsername,
	}
	if target.AccountNo == "" {
		a := accounts.
			NewAccountDomain().
			GetAccountByUserIdAndType(target.UserId, services.EnvelopeAccountType)
		target.AccountNo = a.AccountNo
	}
	transferDTO := services.AccountTransferDTO{
		TradeNo:     dto.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		// 本次抢到红包金额
		Amount:     domain.itemDomain.RedEnvelopeItem.Amount,
		ChangeType: services.EnvelopeOutgoing,
		ChangeFlag: services.FlagTransferOut,
		Desc:       "系统红包账户扣减,转入收红包账户,红包编号 " + dto.EnvelopeNo,
	}
	accountDomain := accounts.NewAccountDomain()
	// 系统账户扣减资金
	return accountDomain.TransferWithContextTx(ctx, transferDTO)
}
