package transactions

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/apierr"
	"GoAuth/internal/ports/inbounds"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TxKey struct{}

var TxKeyValue = TxKey{}

type TxRunner struct {
	pool *pgxpool.Pool // Changed from *sql.DB
}

func NewTxRunner(pool *pgxpool.Pool) inbounds.TxRunner {
	return &TxRunner{pool: pool}
}

// WithinTx executes fn inside a transaction using default options
// (serializable isolation, read-write).
func (r *TxRunner) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.WithinTxWithOptions(ctx, inbounds.TxOptions{}, fn)
}

func (r *TxRunner) WithinTxWithOptions(
	ctx context.Context,
	opts inbounds.TxOptions,
	fn func(ctx context.Context) error,
) (err error) {
	if ctx == nil {
		return fail.New(apierr.SYSTransactionNilContext)
	}

	if ctx.Value(TxKeyValue) != nil {
		return fail.New(apierr.DBNestedTransactionNotAllowed)
	}

	pgxOpts := pgx.TxOptions{
		IsoLevel:   opts.Isolation,
		AccessMode: opts.ReadOnly,
	}

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgxOpts)
	if err != nil {
		return fail.New(apierr.DBBeginTransactionFailed).With(err)
	}

	committed := false

	defer func() {
		if p := recover(); p != nil {
			if !committed {
				rbErr := tx.Rollback(ctx)
				if rbErr != nil {
					logs.L().Error("error during tx rollback after panic", zap.Error(rbErr))
				}
			}
			logs.L().Error("transaction function panicked", zap.Any("panic", p))
			err = fail.New(apierr.DBTransactionPanicked).AddMeta("panic", p)
		}
	}()

	ctx = context.WithValue(ctx, TxKeyValue, tx)

	if err = fn(ctx); err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			logs.L().Error("error during tx rollback after usecase error", zap.Error(rbErr))
		}
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		logs.L().Error("error during tx commit", zap.Error(err))
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			logs.L().Error("error during tx rollback after commit failure", zap.Error(rbErr))
		}
		return fail.New(apierr.DBTransactionCommitFailed).With(err)
	}
	committed = true
	return nil
}
