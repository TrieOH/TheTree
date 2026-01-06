package persistence

import (
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/transactions"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type txKey string

var txKeyValue txKey = "tx_key"

type TxRunner struct {
	db *sql.DB
}

func NewTxRunner(db *sql.DB) *TxRunner {
	return &TxRunner{db: db}
}

func (r *TxRunner) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.WithinTxWithOptions(ctx, transactions.TxOptions{}, fn)
}

func (r *TxRunner) WithinTxWithOptions(
	ctx context.Context,
	opts transactions.TxOptions,
	fn func(ctx context.Context) error,
) error {
	sqlOpts := &sql.TxOptions{
		Isolation: opts.Isolation,
		ReadOnly:  opts.ReadOnly,
	}

	tx, err := r.db.BeginTx(ctx, sqlOpts)
	if err != nil {
		return apierr.ErrInternal.WithMsg("cannot begin transaction").WithID(apierr.DBBeginTXFailed).WithCause(err)
	}

	ctx = context.WithValue(ctx, txKeyValue, tx)

	if err := fn(ctx); err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			logs.L().Error("error during tx rollback", zap.Error(txErr))
		}
		return err
	}

	txErr := tx.Commit()
	if txErr != nil {
		logs.L().Error("error during tx commit", zap.Error(txErr))
		return apierr.ErrInternal.WithMsg("transaction commit failed").WithID(apierr.DBCommitTXFailed).WithCause(txErr)
	}
	return nil
}
