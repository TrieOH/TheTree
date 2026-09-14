// Test-only transaction runner stub. RunTx needs a default runner; unit
// tests have no DB, so transactions flatten to direct calls.
package authn

import (
	"context"

	"lib/database"
)

type nopTxRunner struct{}

func (nopTxRunner) WithinTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func (nopTxRunner) WithinTxWithOptions(ctx context.Context, _ database.TxOptions, fn func(context.Context) error) error {
	return fn(ctx)
}

func init() {
	database.SetDefaultRunner(nopTxRunner{})
}
