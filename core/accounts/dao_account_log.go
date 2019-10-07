package accounts

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountLogDao struct {
	runner *dbx.TxRunner
}

// 通过流水编号查询流水记录
func (dao *AccountLogDao) GetOne(logNo string) *AccountLog {
	out := &AccountLog{LogNo: logNo}
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

// 通过交易编号查询流水记录
func (dao *AccountLogDao) GetByTradeNo(tradeNo string) *AccountLog {
	sql := "select * from account_log where trade_no = ?"
	out := &AccountLog{}
	ok, err := dao.runner.Get(out, sql, tradeNo)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return out
}

// 流水记录的写入
func (dao *AccountLogDao) Insert(data *AccountLog) (id int64, err error) {
	result, err := dao.runner.Insert(data)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return result.LastInsertId()
}
