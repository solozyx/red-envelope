package envelopes

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type RedEnvelopeItemDao struct {
	runner *dbx.TxRunner
}

func (dao *RedEnvelopeItemDao) GetOne(itemNo string) *RedEnvelopeItem {
	out := &RedEnvelopeItem{ItemNo: itemNo}
	ok, err := dao.runner.GetOne(out)
	if !ok || err != nil {
		logrus.Error(err)
		return nil
	}
	return out
}

func (dao *RedEnvelopeItemDao) Insert(data *RedEnvelopeItem) (int64, error) {
	rs, err := dao.runner.Insert(data)
	if err != nil {
		return 0, nil
	}
	return rs.LastInsertId()
}

func (dao *RedEnvelopeItemDao) FindItems(envelopeNo string) []*RedEnvelopeItem {
	items := make([]*RedEnvelopeItem, 0)
	sql := "select * from red_envelope_item where envelope_no = ?"
	err := dao.runner.Find(&items, sql, envelopeNo)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return items
}

func (dao *RedEnvelopeItemDao) ListReceivedItems(userId string, offset, limit int) []*RedEnvelopeItem {
	items := make([]*RedEnvelopeItem, 0)
	sqlQuery := "select * from red_envelope_item where recv_user_id = ? order by created_at desc limit ?,?"
	err := dao.runner.Find(&items, sqlQuery, userId, offset, limit)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return items
}

func (dao *RedEnvelopeItemDao) GetByUser(envelopeNo, userId string) *RedEnvelopeItem {
	item := RedEnvelopeItem{}
	sql := "select * from red_envelope_item where envelope_no=? and recv_user_id=?"
	ok, err := dao.runner.Get(&item, sql, envelopeNo, userId)
	if !ok {
		return nil
	}
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return &item
}
