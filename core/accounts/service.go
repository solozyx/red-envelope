package accounts

import (
	"errors"
	"fmt"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
)

var _ services.AccountService = new(accountService)

var one sync.Once

func init() {
	one.Do(func() {
		services.IAccountService = new(accountService)
	})
}

type accountService struct{}

// 创建账户
func (s *accountService) CreateAccount(dto services.AccountCreatedDTO) (*services.AccountDTO, error) {
	domain := accountDomain{}
	// 验证输入参数
	err := base.ValidateStruct(&dto)
	if err != nil {
		return nil, err
	}
	// 验证账户是否已经存在
	acc := domain.GetAccountByUserIdAndType(dto.UserId, services.EnvelopeAccountType)
	if acc != nil {
		return acc, errors.New(fmt.Sprintf("用户的该类型账户已经存在，username=%s[%s],账户类型:%d",
			acc.Username, acc.UserId, acc.AccountType))
	}
	// 执行账户创建的业务
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	// services.AccountCreatedDTO --> services.AccountDTO
	account := services.AccountDTO{
		UserId:       dto.UserId,
		Username:     dto.Username,
		AccountName:  dto.AccountName,
		AccountType:  dto.AccountType,
		CurrencyCode: dto.CurrencyCode,
		Status:       1,
		Balance:      amount,
	}
	return domain.Create(account)
}

// 转账
func (s *accountService) Transfer(dto services.AccountTransferDTO) (services.TransferredStatus, error) {
	domain := accountDomain{}
	// 验证dto参数
	err := base.ValidateStruct(&dto)
	if err != nil {
		return services.TransferredStatusFailure, err
	}
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return services.TransferredStatusFailure, err
	}
	dto.Amount = amount
	if dto.ChangeFlag == services.FlagTransferOut {
		if dto.ChangeType > 0 {
			return services.TransferredStatusFailure,
				errors.New("如果changeFlag为支出，那么changeType必须小于0")
		}
	} else {
		if dto.ChangeType < 0 {
			return services.TransferredStatusFailure,
				errors.New("如果changeFlag为收入，那么changeType必须大于0")
		}
	}
	// 执行转账操作
	return domain.Transfer(dto)
}

// 储值
func (s *accountService) StoreValue(dto services.AccountTransferDTO) (services.TransferredStatus, error) {
	// 交易对象 和 交易主体 都是自己 转账的特殊形式
	dto.TradeTarget = dto.TradeBody
	// 入账
	dto.ChangeFlag = services.FlagTransferIn
	// 储值
	dto.ChangeType = services.AccountStoreValue
	// 转账储值
	return s.Transfer(dto)
}

func (s *accountService) GetAccount(accountNo string) *services.AccountDTO {
	domain := accountDomain{}
	return domain.GetAccount(accountNo)
}

// 通过 user_id 查询红包账户
func (s *accountService) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	domain := accountDomain{}
	return domain.GetEnvelopeAccountByUserId(userId)
}
