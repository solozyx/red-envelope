package services

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/solozyx/red-envelope/infra/base"
)

var IAccountService AccountService

// 用于对外暴露资金账户应用服务 唯一的暴露点
func GetAccountService() AccountService {
	base.Check(IAccountService)
	return IAccountService
}

type AccountService interface {
	// 创建账户
	CreateAccount(dto AccountCreatedDTO) (*AccountDTO, error)
	// 转账
	Transfer(dto AccountTransferDTO) (TransferredStatus, error)
	// 储值
	StoreValue(dto AccountTransferDTO) (TransferredStatus, error)
	// 红包账户查询
	GetEnvelopeAccountByUserId(userId string) *AccountDTO
	GetAccount(accountNo string) *AccountDTO
}

// DTO : Data Transfer Object
// 数据传输对象

// 账户创建
type AccountCreatedDTO struct {
	// 用户编号 账户所属用户
	UserId string `json:"userId" validate:"required"`
	// 用户名称
	Username string `json:"username" validate:"required"`
	// 账户名称 用来说明账户的简短描述 比如xxx积分 xxx零钱
	AccountName string `json:"accountName" validate:"required"`
	// 账户类型 用来区分不同类型的账户 积分账户 会员卡账户 钱包账户 红包账户
	AccountType int `json:"accountType"`
	// 货币类型编码 CNY人民币 EUR欧元 USD美元
	CurrencyCode string `json:"currencyCode"`
	// TODO:NOTICE 金额 在Go中 float32 float64 计算时会丢失精度 所以金额用字符串来传递避免丢失精度
	Amount string `json:"amount" validate:"required"`
}

// 账户
type AccountDTO struct {
	// 集成 AccountCreatedDTO 所有属性
	// AccountCreatedDTO

	// 用户编号 账户所属用户
	UserId string `json:"userId" validate:"required"`
	// 用户名称
	Username string `json:"username" validate:"required"`
	// 账户名称 用来说明账户的简短描述 比如xxx积分 xxx零钱
	AccountName string `json:"accountName" validate:"required"`
	// 账户类型 用来区分不同类型的账户 积分账户 会员卡账户 钱包账户 红包账户
	AccountType int `json:"accountType"`
	// 货币类型编码 CNY人民币 EUR欧元 USD美元
	CurrencyCode string `json:"currencyCode"`
	// TODO:NOTICE 金额 在Go中 float32 float64 计算时会丢失精度 所以金额用字符串来传递避免丢失精度
	Amount string `json:"amount" validate:"required"`

	// 账户编号 账户唯一标识
	AccountNo string `json:"accountNo"`
	// 账户创建时间
	CreatedAt time.Time `json:"createdAt"`

	Balance   decimal.Decimal `json:"balance"`   //账户可用余额
	Status    int             `json:"status"`    //账户状态，账户状态：0账户初始化，1启用，2停用
	UpdatedAt time.Time       `json:"updatedAt"` //更新时间
}

// 账户交易参与者 交易主体 交易对方 信息一致
type TradeParticipator struct {
	// 账户编号
	AccountNo string `validate:"required" json:"accountNo"`
	// 用户编号
	UserId string `validate:"required" json:"userId"`
	// 用户名
	Username string `validate:"required" json:"username"`
}

// 账户转账
type AccountTransferDTO struct {
	// 交易订单号
	TradeNo string `validate:"required" json:"tradeNo"`
	// 交易主体
	TradeBody TradeParticipator `validate:"required" json:"tradeBody"`
	// 交易对方
	TradeTarget TradeParticipator `validate:"required" json:"tradeTarget"`
	// 交易金额
	AmountStr string          `validate:"required" json:"amountStr"`
	Amount    decimal.Decimal ``
	// 转账变化类型
	ChangeType ChangeType `validate:"required,numeric" json:"changeType"`
	// 资金交易是否成功的变化标识
	ChangeFlag ChangeFlag `validate:"required,numeric" json:"changeFlag"`
	// 交易描述
	Desc string `json:"desc"`
}

//账户流水
type AccountLogDTO struct {
	LogNo           string          //流水编号 全局不重复字符或数字，唯一性标识
	TradeNo         string          //交易单号 全局不重复字符或数字，唯一性标识
	AccountNo       string          //账户编号 账户ID
	TargetAccountNo string          //账户编号 账户ID
	UserId          string          //用户编号
	Username        string          //用户名称
	TargetUserId    string          //目标用户编号
	TargetUsername  string          //目标用户名称
	Amount          decimal.Decimal //交易金额,该交易涉及的金额
	Balance         decimal.Decimal //交易后余额,该交易后的余额
	ChangeType      ChangeType      //流水交易类型，0 创建账户，>0 为收入类型，<0 为支出类型，自定义
	ChangeFlag      ChangeFlag      //交易变化标识：-1 出账 1为进账，枚举
	Status          int             //交易状态：
	Decs            string          //交易描述
	CreatedAt       time.Time       //创建时间
}
