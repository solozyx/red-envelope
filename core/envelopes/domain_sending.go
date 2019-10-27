package envelopes

import (
	"context"
	"path"

	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/core/accounts"
	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
)

// TODO:NOTICE 发红包业务 保存红包商品 需要在1个数据库事务执行
//  事务逻辑:
//  保存红包商品和红包金额的支付必须要保证 全部成功 或 全部失败 把这2个数据库操作放到同1个数据库事务在事务中保证它们的原子性
// 红包金额支付
// 1.需要红包中间商的红包资金账户 定义在配置文件中 事先初始化到资金账户表中
// 2.从红包发送人的资金账户中扣减红包金额
// 3.将扣减的红包总金额转入红包中间商的红包资金账户
func (domain *goodsDomain) SendOut(dto services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 创建红包商品
	domain.Create(dto)
	// 创建红包活动
	activity = new(services.RedEnvelopeActivity)
	// 红包链接 格式 http://域名/v1/envelope/{id}/link/
	link := base.GetEnvelopeActivityLink()
	domainName := base.GetEnvelopeDomain()
	activity.Link = path.Join(domainName, link, domain.EnvelopeNo)

	accountDomain := accounts.NewAccountDomain()

	err = base.Tx(func(runner *dbx.TxRunner) error {
		// 创建 Context 上下文对象绑定 *dbx.TxRunner 数据库事务对象
		ctx := base.WithValueContext(context.Background(), runner)
		// 1.保存红包商品
		id, err := domain.Save(ctx)
		if id <= 0 || err != nil {
			return err
		}
		// 2.把资金从红包发送人的资金账户里扣除
		// 交易主体 发红包账户
		body := services.TradeParticipator{
			AccountNo: dto.AccountNo,
			UserId:    dto.UserId,
			Username:  dto.Username,
		}
		// 交易对方 系统红包账户
		systemAccount := base.GetSystemAccount()
		target := services.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			UserId:    systemAccount.UserId,
			Username:  systemAccount.Username,
		}
		// 出账 dto
		transfer := services.AccountTransferDTO{
			// 红包id作为交易流水号 红包模块 和 资金模块 桥梁
			TradeNo:     domain.RedEnvelopeGoods.EnvelopeNo,
			TradeBody:   body,
			TradeTarget: target,
			Amount:      domain.RedEnvelopeGoods.Amount,
			ChangeType:  services.EnvelopeOutgoing,
			ChangeFlag:  services.FlagTransferOut,
			Desc:        "发红包人红包金额支付",
		}
		// 发红包
		status, err := accountDomain.TransferWithContextTx(ctx, transfer)
		if status != services.TransferredStatusSuccess {
			return err
		}
		// 3.将扣减红包发送人的红包总金额 转入 红包中间商的红包资金账户
		// 入账 dto
		// TODO ???
		// transfer = services.AccountTransferDTO{
		//	TradeNo:     domain.EnvelopeNo,
		//	TradeBody:   target,
		//	TradeTarget: body,
		//	Amount:      domain.Amount,
		//	ChangeType:  services.EnvelopeIncoming,
		//	ChangeFlag:  services.FlagTransferIn,
		//	Desc:        "红包金额转入",
		//}
		//status, err = accountDomain.TransferWithContextTx(ctx, transfer)
		//if status != services.TransferredStatusSuccess {
		//	return err
		//}
		return nil
	})

	if err != nil {
		return nil, err
	}
	// 扣减金额没有问题 返回红包活动
	activity.RedEnvelopeGoodsDTO = *domain.RedEnvelopeGoods.ToDTO()

	return activity, err
}
