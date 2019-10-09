package envelopes

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"

	"github.com/solozyx/red-envelope/services"
)

// 红包商品 映射 red_envelope_goods 表 dbx的ORM映射tag是 `db`
type RedEnvelopeGoods struct {
	Id               int64                `db:"id,omitempty"`
	EnvelopeNo       string               `db:"envelope_no,unique"`
	EnvelopeType     int                  `db:"envelope_type"`
	Username         sql.NullString       `db:"username"`
	UserId           string               `db:"user_id"`
	Blessing         sql.NullString       `db:"blessing"`
	Amount           decimal.Decimal      `db:"amount"`
	AmountOne        decimal.Decimal      `db:"amount_one"`
	Quantity         int                  `db:"quantity"`
	RemainAmount     decimal.Decimal      `db:"remain_amount"`
	RemainQuantity   int                  `db:"remain_quantity"`
	ExpiredAt        time.Time            `db:"expired_at"`
	Status           services.OrderStatus `db:"status"`
	OrderType        services.OrderType   `db:"order_type"`
	PayStatus        services.PayStatus   `db:"pay_status"`
	CreatedAt        time.Time            `db:"created_at,omitempty"`
	UpdatedAt        time.Time            `db:"updated_at,omitempty"`
	OriginEnvelopeNo string               `db:"origin_envelope_no"` // 原关联订单号
}

func (po *RedEnvelopeGoods) ToDTO() *services.RedEnvelopeGoodsDTO {
	return &services.RedEnvelopeGoodsDTO{
		EnvelopeNo:       po.EnvelopeNo,
		EnvelopeType:     po.EnvelopeType,
		Username:         po.Username.String,
		UserId:           po.UserId,
		Blessing:         po.Blessing.String,
		Amount:           po.Amount,
		AmountOne:        po.AmountOne,
		Quantity:         po.Quantity,
		RemainAmount:     po.RemainAmount,
		RemainQuantity:   po.RemainQuantity,
		ExpiredAt:        po.ExpiredAt,
		Status:           int(po.Status),
		OrderType:        po.OrderType,
		PayStatus:        po.PayStatus,
		CreatedAt:        po.CreatedAt,
		UpdatedAt:        po.UpdatedAt,
		AccountNo:        "",
		OriginEnvelopeNo: po.OriginEnvelopeNo,
	}
}

func (po *RedEnvelopeGoods) FromDTO(dto *services.RedEnvelopeGoodsDTO) {
	po.EnvelopeType = dto.EnvelopeType
	po.UserId = dto.UserId
	po.Username = sql.NullString{
		String: dto.Username,
		Valid:  true,
	}
	po.Blessing = sql.NullString{
		String: dto.Blessing,
		Valid:  true,
	}
	po.Amount = dto.Amount
	po.AmountOne = dto.AmountOne
	po.Quantity = dto.Quantity
	po.RemainAmount = dto.RemainAmount
	po.RemainQuantity = dto.RemainQuantity
	po.ExpiredAt = dto.ExpiredAt
	po.Status = services.OrderStatus(dto.Status)
	po.OrderType = dto.OrderType
	po.PayStatus = dto.PayStatus
	po.OriginEnvelopeNo = dto.OriginEnvelopeNo
}
