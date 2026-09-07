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

// WithoutTx returns a context that no longer carries the transaction stored
// under TxKeyValue, if any. Cancellation and all other values are preserved.
//
// A pgx connection can serve only one query at a time: passing a context that
// still holds an open transaction to a goroutine or callback makes its queries
// run on the caller's transaction connection concurrently, which fails with
// pgconn.connLockError ("conn busy") and can corrupt the per-connection
// prepared-statement cache (SQLSTATE 08P01). Detach the tx before handing the
// context to any concurrent work so the queries fall back to the pool.
func WithoutTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, TxKeyValue, nil)
}
