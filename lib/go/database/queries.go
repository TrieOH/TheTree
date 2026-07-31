package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TxQueries[T any] interface {
	WithTx(pgx.Tx) T
}

func Queries[T TxQueries[T]](ctx context.Context, q T) T { //nolint:ireturn
	if tx, ok := ctx.Value(TxKeyValue).(pgx.Tx); ok && tx != nil {
		return q.WithTx(tx)
	}
	return q
}
