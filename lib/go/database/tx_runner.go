package database

import "context"

// TxRunner executes functions within a database transaction.
//
// Implementations are expected to:
//   - Use the database's default isolation, read-write
//   - Reject nested transactions rather than flattening or using save points
//
// A transaction-bound context is passed to fn and must be used by repositories
// to access the active transaction.
//
// There is no package-level default runner: the runner is constructed at
// boot (see PgxTxRunner) and threaded explicitly into the services that
// open transactions, so transaction knowledge stays local to its users and
// no init-order contract exists.
type TxRunner interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// NoTx is the TxRunner for tests that have no database: it calls fn
// directly, so the repo fakes see the caller's context.
type NoTx struct{}

func (NoTx) WithinTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}
