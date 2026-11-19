package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"lib/telemetry"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxKey struct{}

var TxKeyValue = TxKey{}

type PgxTxRunner struct {
	pool *pgxpool.Pool
}

func NewPGXTxRunner(pool *pgxpool.Pool) TxRunner { //nolint:ireturn
	return &PgxTxRunner{pool: pool}
}

// WithinTx executes fn inside a transaction with the database's default
// isolation, read-write.
func (r *PgxTxRunner) WithinTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	if ctx == nil {
		return errors.New("transaction had nil context")
	}

	if ctx.Value(TxKeyValue) != nil {
		return errors.New("nested transactions not allowed")
	}

	var tx pgx.Tx
	tx, err = r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	committed := false

	defer func() {
		if p := recover(); p != nil {
			if !committed {
				rbErr := tx.Rollback(ctx)
				if rbErr != nil {
					telemetry.Log().Error("error during tx rollback after panic", telemetry.Err(rbErr))
				}
			}
			telemetry.Log().Error("transaction function panicked", slog.Any("panic", p))
			err = fmt.Errorf("transaction panicked: %v", p)
		}
	}()

	ctx = context.WithValue(ctx, TxKeyValue, tx)

	err = fn(ctx)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			telemetry.Log().Error("error during tx rollback after usecase error", telemetry.Err(rbErr))
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		telemetry.Log().Error("error during tx commit", telemetry.Err(err))
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			telemetry.Log().Error("error during tx rollback after commit failure", telemetry.Err(rbErr))
		}
		return fmt.Errorf("commit transaction: %w", err)
	}
	committed = true
	return nil
}
