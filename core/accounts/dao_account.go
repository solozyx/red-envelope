package accounts

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

// 查询数据库持久化对象的单实例
// 数据库访问层的每个方法 代表1个原子性操作
// 不建议把多个数据库操作写到1个方法中
// 数据库操作的事务放到各个方法外部,多个方法对应的数据库操作组成1个事务
type AccountDao struct {
	// 执行数据库事务操作
	runner *dbx.TxRunner
}

// 获取一行数据
func (dao *AccountDao) GetOne(accountNo string) *Account {
	// GetOne 方法必须传入1个唯一索引字段 主键 唯一索引 都可以 account_no 是唯一索引字段
	a := &Account{AccountNo: accountNo}
	ok, err := dao.runner.GetOne(a)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

// 资金账户Account 依赖用户User 1个user可以有多个account通过账户类型区分
// 通过用户Id和账户类型来查询账户信息
func (dao *AccountDao) GetByUserId(userId string, accountType int) *Account {
	a := &Account{}
	sql := `select * from account where user_id=? and account_type=?`
	ok, err := dao.runner.Get(a, sql, userId, accountType)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

// 账户数据插入 返回账户id
func (dao *AccountDao) Insert(data *Account) (int64, error) {
	result, err := dao.runner.Insert(data)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return result.LastInsertId()
}

// TODO:NOTICE 账户余额更新 [乐观锁] 机制
// 账户余额更新 amount > 0 收入 amount < 0 扣减
// 入参 account 变化的金额 负为扣减 正为增加
// 返回受影响行数 成功返回1 失败返回<=0的值
func (dao *AccountDao) UpdateBalance(accountNo string, amount decimal.Decimal) (rows int64, err error) {
	// 该SQL语句使用 [乐观锁] 概念 在 where 条件加上限制 避免金额扣减为负数
	// and balance>=-1*CAST(? as DECIMAL(30,6))
	sql := "update account " +
		" set balance=balance+CAST(? as DECIMAL(30,6)) " +
		" where account_no=? " +
		" and balance>=-1*CAST(? as DECIMAL(30,6))"
	// amount 金额 decimal类型 不能直接被Go的database/sql包识别 转为字符串
	// balance=balance+CAST(? as DECIMAL(30,6)) 长度30位 小数6位
	// CAST函数把字符串在SQL中转为decimal类型 就可以和 balance进行加减
	rs, err := dao.runner.Exec(sql, amount.String(), accountNo, amount.String())
	if err != nil {
		return 0, err
	}
	// 返回受影响行数 成功返回1 失败返回<=0的值
	return rs.RowsAffected()
}

// 账户状态更新 账户数据不做物理删除通过status字段标识账户的 启用 停用 状态
func (dao *AccountDao) UpdateStatus(accountNo string, status int) (rows int64, err error) {
	sql := "update account set status = ? where account_no = ?"
	rs, err := dao.runner.Exec(sql, status, accountNo)
	if err != nil {
		logrus.Error(err)
		return 0, nil
	}
	return rs.RowsAffected()
}
