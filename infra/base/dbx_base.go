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

// 将runner绑定到上下文
func WithValueContext(parent context.Context, runner *dbx.TxRunner) context.Context {
	return context.WithValue(parent, TX, runner)
}

func ExecuteContext(ctx context.Context, fn func(runner *dbx.TxRunner) error) error {
	tx, ok := ctx.Value(TX).(*dbx.TxRunner)
	if !ok || tx == nil {
		logrus.Panic("是否在事务函数块中使用？")
	}
	return fn(tx)
}
