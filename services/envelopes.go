package services

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/solozyx/red-envelope/infra/base"
)

var IRedEnvelopeService RedEnvelopeService

func GetRedEnvelopeService() RedEnvelopeService {
	base.Check(IRedEnvelopeService)
	return IRedEnvelopeService
}

// 红包服务接口
type RedEnvelopeService interface {
	// 发红包
	SendOut(RedEnvelopeSendingDTO) (*RedEnvelopeActivity, error)
	// 收红包 返回订单详情信息
	Receive(RedEnvelopeReceiveDTO) (*RedEnvelopeItemDTO, error)
	// 退款 返回红包商品信息
	Refund(envelopeNo string) (order *RedEnvelopeGoodsDTO)
	// 查询红包订单
	Get(envelopeNo string) (order *RedEnvelopeGoodsDTO)
	// 查询本人发送的红包列表
	ListSent(string, int, int) []*RedEnvelopeGoodsDTO
	ListReceived(userId string, page, size int) []*RedEnvelopeItemDTO
	ListItems(envelopeNo string) []*RedEnvelopeItemDTO
	ListReceivable(int, int) []*RedEnvelopeGoodsDTO
}

// 发红包
type RedEnvelopeSendingDTO struct {
	EnvelopeType int    `json:"envelopeType" validate:"required"`
	Username     string `json:"username" validate:"required"`
	UserId       string `json:"userId" validate:"required"`
	Blessing     string `json:"blessing"`
	// 根据红包类型 EnvelopeType 区分 普通红包指单个红包金额 碰运气红包指红包总金额
	Amount   decimal.Decimal `json:"amount" validate:"required,numeric"`
	Quantity int             `json:"quantity" validate:"required,numeric"`
}

func (dto *RedEnvelopeSendingDTO) ToGoods() *RedEnvelopeGoodsDTO {
	return &RedEnvelopeGoodsDTO{
		EnvelopeType: dto.EnvelopeType,
		Username:     dto.Username,
		UserId:       dto.UserId,
		Blessing:     dto.Blessing,
		Amount:       dto.Amount,
		Quantity:     dto.Quantity,
	}
}

// 收红包
type RedEnvelopeReceiveDTO struct {
	EnvelopeNo   string `json:"envelopeNo" validate:"required"`
	RecvUsername string `json:"recvUsername" validate:"required"`
	RecvUserId   string `json:"recvUserId" validate:"required"`
	// 内部通过 RecvUserId 查询出账号
	AccountNo string `json:"accountNo"`
}

type RedEnvelopeActivity struct {
	// 红包商品
	RedEnvelopeGoodsDTO
	// 活动链接 动态生成全局唯一 发给收红包的人群
	Link string `json:"link"`
}

func (this *RedEnvelopeActivity) CopyTo(target *RedEnvelopeActivity) {
	target.Link = this.Link
	target.EnvelopeNo = this.EnvelopeNo
	target.EnvelopeType = this.EnvelopeType
	target.Username = this.Username
	target.UserId = this.UserId
	target.Blessing = this.Blessing
	target.Amount = this.Amount
	target.AmountOne = this.AmountOne
	target.Quantity = this.Quantity
	target.RemainAmount = this.RemainAmount
	target.RemainQuantity = this.RemainQuantity
	target.ExpiredAt = this.ExpiredAt
	target.Status = this.Status
	target.OrderType = this.OrderType
	target.PayStatus = this.PayStatus
	target.CreatedAt = this.CreatedAt
	target.UpdatedAt = this.UpdatedAt
}

// 红包商品
type RedEnvelopeGoodsDTO struct {
	EnvelopeNo     string          `json:"envelopeNo"`
	EnvelopeType   int             `json:"envelopeType" validate:"required,numeric"`
	Username       string          `json:"username" validate:"required"`
	UserId         string          `json:"userId" validate:"required"`
	Blessing       string          `json:"blessing"`
	Amount         decimal.Decimal `json:"amount" validate:"required,numeric"`
	AmountOne      decimal.Decimal `json:"amountOne"`
	Quantity       int             `json:"quantity" validate:"required,numeric"`
	RemainAmount   decimal.Decimal `json:"remainAmount"`
	RemainQuantity int             `json:"remainQuantity"`
	ExpiredAt      time.Time       `json:"expiredAt"`
	Status         int             `json:"status"`
	OrderType      OrderType       `json:"orderType"`
	PayStatus      PayStatus       `json:"payStatus"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	AccountNo      string          `json:"accountNo"`
	// 原关联订单号
	OriginEnvelopeNo string `json:"originEnvelopeNo"`
}

// 红包详情
type RedEnvelopeItemDTO struct {
	ItemNo       int64           `json:"itemNo"`       // 红包订单详情编号
	EnvelopeNo   string          `json:"envelopeNo"`   // 红包编号
	RecvUsername string          `json:"recvUsername"` // 接收者用户名
	RecvUserId   string          `json:"recvUserId"`   // 接收者用户id
	Amount       decimal.Decimal `json:"amount"`       // 收到金额
	Quantity     int             `json:"quantity"`     // 收到数量
	RemainAmount decimal.Decimal `json:"remainAmount"` // 剩余金额
	AccountNo    string          `json:"accountNo"`    // 红包接收者账户ID
	PayStatus    int             `json:"payStatus"`    // 支付状态
	CreatedAt    time.Time       `json:"createdAt"`    // 创建时间
	UpdatedAt    time.Time       `json:"updatedAt"`    // 修改时间
	Desc         string          `json:"desc"`
	IsLuckiest   bool            `json:"isLuckiest"` // 是否是最幸运的
}

func (item *RedEnvelopeItemDTO) CopyTo(target *RedEnvelopeItemDTO) {
	target.ItemNo = item.ItemNo
	target.EnvelopeNo = item.EnvelopeNo
	target.RecvUsername = item.RecvUsername
	target.RecvUserId = item.RecvUserId
	target.Amount = item.Amount
	target.Quantity = item.Quantity
	target.RemainAmount = item.RemainAmount
	target.AccountNo = item.AccountNo
	target.PayStatus = item.PayStatus
	target.CreatedAt = item.CreatedAt
	target.UpdatedAt = item.UpdatedAt
}
