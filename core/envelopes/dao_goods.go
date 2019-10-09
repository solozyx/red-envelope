package envelopes

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"

	"github.com/solozyx/red-envelope/services"
)

type RedEnvelopeGoodsDao struct {
	runner *dbx.TxRunner
}

// 插入 返回红包商品id
func (dao *RedEnvelopeGoodsDao) Insert(po *RedEnvelopeGoods) (int64, error) {
	rs, err := dao.runner.Insert(po)
	if err != nil {
		logrus.Error(err)
		return 0, nil
	}
	return rs.LastInsertId()
}

// 根据红包编号查询红包商品
func (dao *RedEnvelopeGoodsDao) GetOne(envelopeNo string) *RedEnvelopeGoods {
	var out = &RedEnvelopeGoods{EnvelopeNo: envelopeNo}
	ok, err := dao.runner.GetOne(out)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return out
}

// TODO:NOTICE 更新红包余额和数量 [乐观锁] unsigned 标识 无符号整型 从数据库层面避免写入负数
//  `remain_amount` decimal(30,6) unsigned not null default '0.000000' comment '红包剩余金额',
//  `remain_quantity` int(10) unsigned not null comment '红包剩余数量',
//  [乐观锁] 核心部分SQL " and remain_quantity > 0 and remain_amount >= CAST(? as DECIMAL(30,6)) "
//  不再使用事务行锁来更新红包余额和数量,避免锁竞争,避免死锁
//  改用乐观锁来保证更新操作的安全性 避免 红包金额 剩余数量 负库存问题
//  通过where子句来判断红包剩余金额和数量来解决2个问题
//  1.负库存问题 避免红包金额和数量不足时进行扣减
//  2.减少实际的库存更新 当红包剩余金额和数量不足时不去更新数据库不操作磁盘 过滤掉无效的更新 提高总体性能
//  在100人的微信群,发个红包,数量10个,如果100人都抢,只有10个人能抢到,过滤掉无效的90次更新数据库操作
//  整体性能提升 30%
// 返回 影响行数
func (dao *RedEnvelopeGoodsDao) UpdateBalance(envelopeNo string, amount decimal.Decimal) (int64, error) {
	// CAST 函数 CAST(? as DECIMAL(30,6)) string --> decimal
	sql := "update red_envelope_goods " +
		" set remain_amount = remain_amount - CAST(? as DECIMAL(30,6)), " +
		" remain_quantity = remain_quantity - 1 " +
		" where envelope_no = ? " +
		" and remain_quantity > 0 " +
		" and remain_amount >= CAST(? as DECIMAL(30,6)) "
	rs, err := dao.runner.Exec(sql, amount.String(), envelopeNo, amount.String())
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	// 实际更新 0 行 就是剩余红包不足
	return rs.RowsAffected()
}

// 更新订单状态
func (dao *RedEnvelopeGoodsDao) UpdateOrderStatus(envelopeNo string, status services.OrderStatus) (int64, error) {
	sql := " update red_envelope_goods set status=? where envelope_no=?"
	rs, err := dao.runner.Exec(sql, status, envelopeNo)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return rs.RowsAffected()
}

// 过期 把所有的过期红包都查询出来 分页 limit offset size
func (dao *RedEnvelopeGoodsDao) FindExpired(offset, size int) []RedEnvelopeGoods {
	var goodsList []RedEnvelopeGoods
	now := time.Now()
	sql := "select * from red_envelope_goods " +
		" where remain_quantity>0 " +
		" and expired_at<? and (status<4 or status>5) " +
		" limit ?,?"
	err := dao.runner.Find(&goodsList, sql, now, offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goodsList
}

func (dao *RedEnvelopeGoodsDao) Find(po *RedEnvelopeGoods, offset, limit int) []RedEnvelopeGoods {
	var redEnvelopeGoodss []RedEnvelopeGoods
	err := dao.runner.FindExample(po, &redEnvelopeGoodss)
	if err != nil {
		logrus.Error(err)
	}
	return redEnvelopeGoodss
}

func (dao *RedEnvelopeGoodsDao) FindByUser(userId string, offset, limit int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods

	sql := " select * from red_envelope_goods " +
		" where  user_id=?  order by created_at desc limit ?,?"
	err := dao.runner.Find(&goods, sql, userId, offset, limit)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}

func (dao *RedEnvelopeGoodsDao) ListReceivable(offset, size int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods
	now := time.Now()
	sql := " select * from red_envelope_goods " +
		" where  remain_quantity>0  and expired_at>? order by created_at desc limit ?,?"
	err := dao.runner.Find(&goods, sql, now, offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}
