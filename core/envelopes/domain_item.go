package envelopes

import (
	"context"

	"github.com/segmentio/ksuid"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/infra/base"
	"github.com/solozyx/red-envelope/services"
)

type itemDomain struct {
	RedEnvelopeItem
}

// 生成 itemNo
func (domain *itemDomain) createItemNo() {
	domain.ItemNo = ksuid.New().Next().String()
}

// 创建 item
func (domain *itemDomain) Create(dto services.RedEnvelopeItemDTO) {
	domain.RedEnvelopeItem.FromDTO(&dto)
	// sql.NullString 要把该字段写入数据库 .Valid = true 这样 ORM框架才能识别
	domain.RecvUsername.Valid = true
	domain.createItemNo()
}

// 保存 item
func (domain *itemDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeItemDao{runner: runner}
		id, err = dao.Insert(&domain.RedEnvelopeItem)
		return err
	})
	return id, err
}

// 通过 itemNo 查询抢红包明细数据
func (domain *itemDomain) GetOne(ctx context.Context, itemNo string) (dto *services.RedEnvelopeItemDTO) {
	err := base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeItemDao{runner: runner}
		po := dao.GetOne(itemNo)
		if po != nil {
			dto = po.ToDTO()
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dto
}

// 通过 envelopeNo 查询已抢红包列表
func (domain *itemDomain) FindItems(envelopeNo string) (itemDTOs []*services.RedEnvelopeItemDTO) {
	var items []*RedEnvelopeItem
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeItemDao{runner: runner}
		items = dao.FindItems(envelopeNo)
		return nil
	})
	if err != nil {
		return nil
	}
	itemDTOs = make([]*services.RedEnvelopeItemDTO, 0)
	for _, po := range items {
		itemDTOs = append(itemDTOs, po.ToDTO())
	}
	return itemDTOs
}
