package database

import (
	"context"
	"univents/internal/plataform/telemetry"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TxKey struct{}

var TxKeyValue = TxKey{}

type PgxTxRunner struct {
	pool *pgxpool.Pool // Changed from *sql.DB
}

func NewPGXTxRunner(pool *pgxpool.Pool) TxRunner {
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
		return fail.New(errx.SYSTransactionNilContext)
	}

	if ctx.Value(TxKeyValue) != nil {
		return fail.New(errx.DBNestedTransactionNotAllowed)
	}

	pgxOpts := pgx.TxOptions{
		IsoLevel:   opts.Isolation,
		AccessMode: opts.ReadOnly,
	}

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgxOpts)
	if err != nil {
		return fail.New(errx.DBBeginTransactionFailed).With(err)
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
			err = fail.New(errx.DBTransactionPanicked).AddMeta("panic", p)
		}
	}()

	ctx = context.WithValue(ctx, TxKeyValue, tx)

	if err = fn(ctx); err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			telemetry.Log().Error("error during tx rollback after usecase error", zap.Error(rbErr))
		}
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		telemetry.Log().Error("error during tx commit", zap.Error(err))
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			telemetry.Log().Error("error during tx rollback after commit failure", zap.Error(rbErr))
		}
		return fail.New(errx.DBTransactionCommitFailed).With(err)
	}
	committed = true
	return nil
}
