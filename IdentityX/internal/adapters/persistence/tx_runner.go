package persistence

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/transactions"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type txKey struct{}

var txKeyValue = txKey{}

type TxRunner struct {
	db *sql.DB
}

func NewTxRunner(db *sql.DB) *TxRunner {
	return &TxRunner{db: db}
}

// WithinTx executes fn inside a transaction using default options
// (database default isolation, read-write).
//
// Nested calls are not supported and will return an error.
func (r *TxRunner) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.WithinTxWithOptions(ctx, transactions.TxOptions{}, fn)
}

func (r *TxRunner) WithinTxWithOptions(
	ctx context.Context,
	opts transactions.TxOptions,
	fn func(ctx context.Context) error,
) (err error) {
	if ctx == nil {
		return apierr.ErrInternal.WithMsg("cannot create transactions with a nil context").WithID(apierr.SystemTransactionWithNoContext)
	}

	if ctx.Value(txKeyValue) != nil {
		// Nested transactions are explicitly not supported to avoid implicit
		// transaction reuse or savepoint semantics.
		return apierr.ErrInternal.WithMsg("nested transactions are not supported").WithID(apierr.DBNestedTXNotAllowed)
	}

	sqlOpts := &sql.TxOptions{
		Isolation: opts.Isolation,
		ReadOnly:  opts.ReadOnly,
	}

	var tx *sql.Tx
	tx, err = r.db.BeginTx(ctx, sqlOpts)
	if err != nil {
		return apierr.ErrInternal.WithMsg("cannot begin transaction").WithID(apierr.DBBeginTXFailed).WithCause(err)
	}

	committed := false

	defer func() {
		if p := recover(); p != nil {
			if !committed {
				rbErr := tx.Rollback()
				if rbErr != nil {
					logs.L().Error("error during tx rollback after panic", zap.Error(rbErr))
				}
			}
			logs.L().Error("transaction function panicked", zap.Any("panic", p))
			err = apierr.ErrInternal.
				WithMsg("transaction function panicked").
				WithID(apierr.SystemInternalError)
		}
	}()

	ctx = context.WithValue(ctx, txKeyValue, tx)

	if err = fn(ctx); err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			logs.L().Error("error during tx rollback after usecase error", zap.Error(rbErr))
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		logs.L().Error("error during tx commit", zap.Error(err))
		if rbErr := tx.Rollback(); rbErr != nil {
			logs.L().Error("error during tx rollback after commit failure", zap.Error(rbErr))
		}
		return apierr.ErrInternal.
			WithMsg("transaction commit failed").
			WithID(apierr.DBCommitTXFailed).
			WithCause(err)
	}
	committed = true
	return nil
}
