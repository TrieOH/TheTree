package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
)

type TxQueries[T any] interface {
	WithTx(pgx.Tx) T
}

func Queries[T TxQueries[T]](ctx context.Context, q T) T {
	if tx, ok := ctx.Value(TxKeyValue).(pgx.Tx); ok && tx != nil {
		return q.WithTx(tx)
	}

	return q
}

func Span(ctx context.Context, tracer trace.Tracer, op string) (context.Context, trace.Span) {
	return tracer.Start(ctx, "FormsRepo."+op)
}
