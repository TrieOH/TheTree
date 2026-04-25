package database

import (
	"context"

	"github.com/MintzyG/FastUtilitiesNet"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TxKey struct{}

var TxKeyValue = TxKey{}

type PgxTxRunner struct {
	logger *zap.Logger
	pool   *pgxpool.Pool // Changed from *sql.DB
}

func NewPGXTxRunner(pool *pgxpool.Pool, logger *zap.Logger) TxRunner {
	return &PgxTxRunner{pool: pool, logger: logger}
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
		return fun.NewError("transaction had nil context").Internal()
	}

	if ctx.Value(TxKeyValue) != nil {
		return fun.NewError("nested transactions not allowed").Internal()
	}

	pgxOpts := pgx.TxOptions{
		IsoLevel:   opts.Isolation,
		AccessMode: opts.ReadOnly,
	}

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgxOpts)
	if err != nil {
		return fun.NewErrorf("error beginning transaction: %s", err.Error()).Internal()
	}

	committed := false

	defer func() {
		if p := recover(); p != nil {
			if !committed {
				rbErr := tx.Rollback(ctx)
				if rbErr != nil {
					r.logger.Error("error during tx rollback after panic", zap.Error(rbErr))
				}
			}
			r.logger.Error("transaction function panicked", zap.Any("panic", p))
			err = fun.NewErrorf("transaction panicked").Internal()
		}
	}()

	ctx = context.WithValue(ctx, TxKeyValue, tx)

	if err = fn(ctx); err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			r.logger.Error("error during tx rollback after usecase error", zap.Error(rbErr))
		}
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		r.logger.Error("error during tx commit", zap.Error(err))
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			r.logger.Error("error during tx rollback after commit failure", zap.Error(rbErr))
		}
		return fun.NewErrorf("error commiting transaction: %s", err.Error()).Internal()
	}
	committed = true
	return nil
}
