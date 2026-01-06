package transactions

import (
	"context"
	"database/sql"
)

type TxOptions struct {
	Isolation sql.IsolationLevel
	ReadOnly  bool
}

type TxRunner interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
	WithinTxWithOptions(ctx context.Context, opts TxOptions, fn func(ctx context.Context) error) error
}
