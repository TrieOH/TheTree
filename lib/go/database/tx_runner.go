package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// TxOptions defines transaction behavior.
// Zero values result in explicit default options being passed:
//   - Isolation: sql.LevelDefault (driver-defined default isolation)
//   - ReadOnly: false
type TxOptions struct {
	Isolation pgx.TxIsoLevel
	ReadOnly  pgx.TxAccessMode
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

var defaultRunner TxRunner

// SetDefaultRunner sets the package-level transaction runner. Call once at startup.
func SetDefaultRunner(r TxRunner) { defaultRunner = r }

// RunTx executes fn within a transaction using the default runner.
func RunTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return defaultRunner.WithinTx(ctx, fn)
}

// RunTxWithOptions executes fn within a transaction with the given options.
func RunTxWithOptions(ctx context.Context, opts TxOptions, fn func(ctx context.Context) error) error {
	return defaultRunner.WithinTxWithOptions(ctx, opts, fn)
}
