// Test-only transaction runner stub: unit tests have no DB, so the
// transactions multi-step writes run in flatten to direct calls.
package oauth_providers

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
