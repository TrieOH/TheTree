package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/ports/outbounds"
	"context"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type tokenReuseListRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *tokenReuseListRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbounds.TokenReuseListRepository = (*tokenReuseListRepo)(nil)

func NewTokenReuseListRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbounds.TokenReuseListRepository {
	return &tokenReuseListRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func (repo *tokenReuseListRepo) Append(ctx context.Context, jit, userID uuid.UUID, expiresAt time.Time) error {
	ctx, span := repo.tracer.Start(ctx, "TokenReuseListRepo.Append")
	defer span.End()

	err := repo.queries(ctx).TokenReuseListAppend(ctx, sqlc.TokenReuseListAppendParams{
		Jit:       jit,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *tokenReuseListRepo) Exists(ctx context.Context, jit, userID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "TokenReuseListRepo.Exists")
	defer span.End()

	exists, err := repo.queries(ctx).TokenReuseListExists(ctx, sqlc.TokenReuseListExistsParams{
		Jit:    jit,
		UserID: userID,
	})
	if err != nil {
		return false, fail.From(err).RecordCtx(ctx)
	}
	return exists, nil
}

func (repo *tokenReuseListRepo) ClearExpired(ctx context.Context) error {
	ctx, span := repo.tracer.Start(ctx, "TokenReuseListRepo.ClearExpired")
	defer span.End()

	err := repo.queries(ctx).DeleteExpiredTokenReuseListEntries(ctx)
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}
	return nil
}
