package base

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

const TX = "tx"

type BaseDao struct {
	Tx *sql.Tx
}

func (d *BaseDao) SetTx(tx *sql.Tx) {
	d.Tx = tx
}

type TxFunc func(runner *dbx.TxRunner) error

// 事务执行
func Tx(fn TxFunc) error {
	return TxContext(context.Background(), fn)
}

// 事务执行
func TxContext(ctx context.Context, fn TxFunc) error {
	return DbxDatabase().Tx(fn)
}

// 在context上下文对象中传递 *TxRunner对象
// 实现跨方法执行事务 在不同方法使用同1个 *TxRunner事务对象 在同1个*TxRunner事务对象执行不同逻辑
// 和 ExecuteContext 配合
func WithValueContext(parent context.Context, runner *dbx.TxRunner) context.Context {
	return context.WithValue(parent, TX, runner)
}

// 该方法没有真正开启1个事务,从 ctx 获取 *TxRunner 事务对象
func ExecuteContext(ctx context.Context, fn TxFunc) error {
	tx, ok := ctx.Value(TX).(*dbx.TxRunner)
	if !ok || tx == nil {
		logrus.Panic("是否在事务函数块中使用?")
	}
	return fn(tx)
}
