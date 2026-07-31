package database

import (
	"context"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TxKey struct{}

var TxKeyValue = TxKey{}

type PgxTxRunner struct {
	pool *pgxpool.Pool // Changed from *sql.DB
}

func NewPGXTxRunner(pool *pgxpool.Pool) TxRunner { //nolint:ireturn
	return &PgxTxRunner{pool: pool}
}

// WithinTx executes fn inside a transaction using default options
// (serializable isolation, read-write).
func (r *PgxTxRunner) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.WithinTxWithOptions(ctx, TxOptions{}, fn)
}

func (r *PgxTxRunner) WithinTxWithOptions(
	ctx context.Context,
	opts TxOptions,
	fn func(ctx context.Context) error,
) (err error) {
	if ctx == nil {
		return fun.ErrInternal("transaction had nil context")
	}

	if ctx.Value(TxKeyValue) != nil {
		return fun.ErrInternal("nested transactions not allowed")
	}

	pgxOpts := pgx.TxOptions{
		IsoLevel:   opts.Isolation,
		AccessMode: opts.ReadOnly,
	}

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgxOpts)
	if err != nil {
		return fun.Errf("error beginning transaction: %s", err.Error()).Internal()
	}

	committed := false

	defer func() {
		if p := recover(); p != nil {
			if !committed {
				rbErr := tx.Rollback(ctx)
				if rbErr != nil {
					telemetry.Log().Error("error during tx rollback after panic", zap.Error(rbErr))
				}
			}
			telemetry.Log().Error("transaction function panicked", zap.Any("panic", p))
			err = fun.ErrInternal("transaction panicked")
		}
	}()

	ctx = context.WithValue(ctx, TxKeyValue, tx)

	err = fn(ctx)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			telemetry.Log().Error("error during tx rollback after usecase error", zap.Error(rbErr))
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		telemetry.Log().Error("error during tx commit", zap.Error(err))
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			telemetry.Log().Error("error during tx rollback after commit failure", zap.Error(rbErr))
		}
		return fun.Errf("error commiting transaction: %s", err.Error()).Internal()
	}
	committed = true
	return nil
}
