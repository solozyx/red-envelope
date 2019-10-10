package envelopes

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"

	"github.com/solozyx/red-envelope/services"
)

type RedEnvelopeItem struct {
	Id           int64           `db:"id,omitempty"`   // 自增id
	ItemNo       string          `db:"item_no,unique"` // 红包订单详情编号
	EnvelopeNo   string          `db:"envelope_no"`    // 红包编号
	RecvUsername sql.NullString  `db:"recv_username"`  // 接收者用户名
	RecvUserId   string          `db:"recv_user_id"`   // 接收者用户id
	Amount       decimal.Decimal `db:"amount"`         // 收到金额
	Quantity     int             `db:"quantity"`       // 收到数量
	RemainAmount decimal.Decimal `db:"remain_amount"`  // 剩余金额
	AccountNo    string          `db:"account_no"`     // 红包接收者账户ID
	PayStatus    int             `db:"pay_status"`     // 支付状态
	Desc         string          `db:"desc"`
	CreatedAt    time.Time       `db:"created_at,omitempty"` // 创建时间
	UpdatedAt    time.Time       `db:"updated_at,omitempty"` // 修改时间
}

func (po *RedEnvelopeItem) ToDTO() *services.RedEnvelopeItemDTO {
	return &services.RedEnvelopeItemDTO{
		ItemNo:       po.ItemNo,
		EnvelopeNo:   po.EnvelopeNo,
		RecvUsername: po.RecvUsername.String,
		RecvUserId:   po.RecvUserId,
		Amount:       po.Amount,
		Quantity:     po.Quantity,
		RemainAmount: po.RemainAmount,
		AccountNo:    po.AccountNo,
		PayStatus:    po.PayStatus,
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
		Desc:         po.Desc,
	}
}

func (po *RedEnvelopeItem) FromDTO(dto *services.RedEnvelopeItemDTO) {
	po.ItemNo = dto.ItemNo
	po.EnvelopeNo = dto.EnvelopeNo
	po.RecvUsername = sql.NullString{
		String: dto.RecvUsername,
		Valid:  true,
	}
	po.RecvUserId = dto.RecvUserId
	po.Amount = dto.Amount
	po.Quantity = dto.Quantity
	po.RemainAmount = dto.RemainAmount
	po.AccountNo = dto.AccountNo
	po.PayStatus = dto.PayStatus
	po.Desc = dto.Desc
}
