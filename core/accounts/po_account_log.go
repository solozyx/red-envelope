package accounts

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/solozyx/red-envelope/services"
)

// 账户流水表持久化对象
type AccountLog struct {
	Id              int64               `db:"id,omitempty"`
	LogNo           string              `db:"log_no,unique"`
	TradeNo         string              `db:"trade_no"`
	AccountNo       string              `db:"account_no"`
	UserId          string              `db:"user_id"`
	Username        string              `db:"username"`
	TargetAccountNo string              `db:"target_account_no"`
	TargetUserId    string              `db:"target_user_id"`
	TargetUsername  string              `db:"target_username"`
	Amount          decimal.Decimal     `db:"amount"`
	Balance         decimal.Decimal     `db:"balance"`
	ChangeType      services.ChangeType `db:"change_type"`
	ChangeFlag      services.ChangeFlag `db:"change_flag"`
	Status          int                 `db:"status"`
	Desc            string              `db:"desc"`
	// 创建时间系统自动生成 该字段无需手动赋值
	CreatedAt time.Time `db:"created_at,omitempty"`
}

func (po *AccountLog) FromTransferDTO(dto *services.AccountTransferDTO) {
	po.TradeNo = dto.TradeNo
	po.AccountNo = dto.TradeBody.AccountNo
	po.UserId = dto.TradeBody.UserId
	po.Username = dto.TradeBody.Username
	po.TargetAccountNo = dto.TradeTarget.AccountNo
	po.TargetUserId = dto.TradeTarget.UserId
	po.TargetUsername = dto.TradeTarget.Username
	po.Amount = dto.Amount
	po.ChangeFlag = dto.ChangeFlag
	po.ChangeType = dto.ChangeType
	po.Desc = dto.Desc
}

func (po *AccountLog) ToDTO() *services.AccountLogDTO {
	return &services.AccountLogDTO{
		LogNo:           po.LogNo,
		TradeNo:         po.TradeNo,
		AccountNo:       po.AccountNo,
		TargetAccountNo: po.TargetAccountNo,
		UserId:          po.UserId,
		Username:        po.Username,
		TargetUserId:    po.TargetUserId,
		TargetUsername:  po.TargetUsername,
		Amount:          po.Amount,
		Balance:         po.Balance,
		ChangeType:      po.ChangeType,
		ChangeFlag:      po.ChangeFlag,
		Status:          po.Status,
		Decs:            po.Desc,
		CreatedAt:       po.CreatedAt,
	}
}
