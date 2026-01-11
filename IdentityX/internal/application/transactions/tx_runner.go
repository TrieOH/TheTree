package transactions

import (
	"context"
	"database/sql"
)

// TxOptions defines transaction behavior.
// Zero values result in explicit default options being passed:
//   - Isolation: sql.LevelDefault (driver-defined default isolation)
//   - ReadOnly: false
type TxOptions struct {
	Isolation sql.IsolationLevel
	ReadOnly  bool
}

// TxRunner executes functions within a database transaction.
//
// Implementations are expected to:
//   - Use database default isolation and read-write mode unless specified
//   - Reject nested transactions rather than flattening or using save points
//
// A transaction-bound context is passed to fn and must be used by repositories
// to access the active transaction.
type TxRunner interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
	WithinTxWithOptions(ctx context.Context, opts TxOptions, fn func(ctx context.Context) error) error
}
