package accounts

import (
	"context"

	"github.com/kataras/iris/core/errors"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
)

// 领域对象 有状态的 每次使用时都要实例化
type accountDomain struct {
	account    Account
	accountLog AccountLog
}

func NewAccountDomain() *accountDomain {
	return new(accountDomain)
}

func (domain *accountDomain) GetAccountNo() string {
	return domain.account.AccountNo
}

// 创建 log_no 逻辑
func (domain *accountDomain) createAccountLogNo() {
	// 暂时采用 ksuid 的 Id 生成策略来创建
	// 后期会优化成可读性比较好的 分布式ID 全局唯一
	domain.accountLog.LogNo = ksuid.New().Next().String()
}

// 生成 account_no 逻辑
func (domain *accountDomain) createAccountNo() {
	domain.account.AccountNo = ksuid.New().Next().String()
}

// 创建流水记录
func (domain *accountDomain) createAccountLog() {
	// 通过 account 来创建流水
	domain.accountLog = AccountLog{}
	domain.createAccountLogNo()
	domain.accountLog.TradeNo = domain.accountLog.LogNo

	// 流水中的交易主体信息
	domain.accountLog.AccountNo = domain.account.AccountNo
	domain.accountLog.UserId = domain.account.UserId
	domain.accountLog.Username = domain.account.Username.String

	// 流水中的交易对象信息
	domain.accountLog.TargetAccountNo = domain.account.AccountNo
	domain.accountLog.TargetUserId = domain.account.UserId
	domain.accountLog.TargetUsername = domain.account.Username.String

	// 金额余额 在账户创建时相同
	domain.accountLog.Amount = domain.account.Balance
	domain.accountLog.Balance = domain.account.Balance

	// 交易变化属性
	domain.accountLog.Desc = "账户创建"
	domain.accountLog.ChangeType = services.AccountCreated
	domain.accountLog.ChangeFlag = services.FlagAccountCreated
}

// 账户创建的业务逻辑代码
func (domain *accountDomain) Create(dto services.AccountDTO) (*services.AccountDTO, error) {
	// 创建账户持久化对象
	domain.account = Account{}
	// DTO 转换为 DAO
	domain.account.FromDTO(&dto)
	domain.createAccountNo()
	// sql.NullString 类型的 Valid = true 才能写入数据库
	domain.account.Username.Valid = true

	// 创建账户流水持久化对象
	domain.createAccountLog()

	// TODO:NOTICE 持久化DAO对象是有状态的 每次操作都要使用新的
	//  持久化账户和流水这2个对象在同1个数据库事务中 整个过程要么全部成功 要么全部失败
	accountDao := AccountDao{}
	accountLogDao := AccountLogDao{}
	var rdto *services.AccountDTO
	// 快捷事务操作函数 base.Tx 在该函数所有的数据库独立操作 被认为构成1个事务
	// 在事务中返回任何非nil的error 事务操作就会失败 在 Tx 中就会把数据库操作回滚
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		accountLogDao.runner = runner
		// 插入账户数据
		id, err := accountDao.Insert(&domain.account)
		if err != nil {
			// 事务回滚
			return err
		}
		if id <= 0 {
			// 插入失败 事务回滚
			return errors.New("创建账户失败")
		}
		// 如果插入成功 再插入流水数据
		id, err = accountLogDao.Insert(&domain.accountLog)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户流水失败")
		}
		domain.account = *accountDao.GetOne(domain.account.AccountNo)
		return nil
	})
	rdto = domain.account.ToDTO()
	return rdto, err
}

// 验证用户该账户是否已经存在
func (domain *accountDomain) GetAccountByUserIdAndType(userId string, aType services.AccountType) *services.AccountDTO {
	var a *Account
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{runner: runner}
		a = accountDao.GetByUserId(userId, int(aType))
		return nil
	})
	if err != nil || a == nil {
		return nil
	} else {
		return a.ToDTO()
	}
}

// 领域对象 转账业务
func (domain *accountDomain) Transfer(dto services.AccountTransferDTO) (status services.TransferredStatus, err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		status, err = domain.transferWithContextTx(ctx, dto)
		return err
	})
	return status, err
}

// 必须在 base.TX 事务块里面运行 不能单独运行
func (domain *accountDomain) transferWithContextTx(ctx context.Context, dto services.AccountTransferDTO) (status services.TransferredStatus, err error) {
	// 如果交易变化是支出类型 修正amount为负值
	var amount = dto.Amount
	if dto.ChangeFlag == services.FlagTransferOut {
		amount = amount.Mul(decimal.NewFromFloat(-1))
	}

	// 创建账户流水记录
	domain.accountLog = AccountLog{}
	domain.accountLog.FromTransferDTO(&dto)
	domain.createAccountLogNo()
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{runner: runner}
		accountLogDao := AccountLogDao{runner: runner}
		// 账户扣减时 检查余额是否足够和更新余额 通过乐观锁验证 余额足够则更新余额
		rows, err := accountDao.UpdateBalance(dto.TradeBody.AccountNo, amount)
		if err != nil {
			status = services.TransferredStatusFailure
			// 返回错误 事务回滚
			return err
		}
		if rows <= 0 && dto.ChangeFlag == services.FlagTransferOut {
			status = services.TransferredStatusSufficientFunds
			return errors.New("余额不足")
		}
		if rows <= 0 && dto.ChangeFlag == services.FlagTransferIn {
			return errors.New("增加余额失败")
		}

		// 转账成功后 写入流水记录
		account := accountDao.GetOne(dto.TradeBody.AccountNo)
		if account == nil {
			return errors.New("查询账户信息出错")
		}
		domain.account = *account
		domain.accountLog.Balance = domain.account.Balance

		// 如果对于交易主体来说，ChangeFlag是资金转出, 则交易目标余额增加
		//if dto.ChangeFlag == services.FlagTransferOut {
		//	rows, err = accountDao.UpdateBalance(dto.TradeTarget.AccountNo, amount.Abs())
		//	if rows < 1 || err != nil {
		//		status = services.TransferredStatusFailure
		//		return errors.New("目标账户余额增加失败")
		//	}
		//}

		id, err := accountLogDao.Insert(&domain.accountLog)
		if err != nil || id <= 0 {
			status = services.TransferredStatusFailure
			// 返回错误 回滚事务
			return errors.New("转账账户流水创建失败")
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
	} else {
		status = services.TransferredStatusSuccess
	}
	return
}

// 根据账户编号来查询账户信息
func (domain *accountDomain) GetAccount(accountNo string) *services.AccountDTO {
	var account *Account
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{runner: runner}
		account = accountDao.GetOne(accountNo)
		return nil
	})

	if err != nil {
		logrus.Error(err)
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

// 根据用户ID来查询红包账户
func (domain *accountDomain) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	var account *Account
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{runner: runner}
		account = accountDao.GetByUserId(userId, int(services.EnvelopeAccountType))
		return nil
	})

	if err != nil {
		logrus.Error(err)
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

// 根据流水Id查询账户流水
func (domain *accountDomain) GetAccountLog(logNo string) *services.AccountLogDTO {
	dao := AccountLogDao{}
	var al *AccountLog
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		al = dao.GetOne(logNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if al == nil {
		return nil
	}
	return al.ToDTO()
}

// 根据交易编号查询账户流水
func (domain *accountDomain) GetAccountLogByTradeNo(tradeNo string) *services.AccountLogDTO {
	dao := AccountLogDao{}
	var al *AccountLog
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		al = dao.GetByTradeNo(tradeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if al == nil {
		return nil
	}
	return al.ToDTO()
}
