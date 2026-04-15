package tokens

import (
	"IdentityX/internal/platform/database"
	sqlc2 "IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/ports"
	"context"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type tokenReuseListRepo struct {
	q      *sqlc2.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *tokenReuseListRepo) queries(ctx context.Context) *sqlc2.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ ports.TokenReuseListRepository = (*tokenReuseListRepo)(nil)

func NewRepo(q *sqlc2.Queries, l *zap.Logger, tracer trace.Tracer) ports.TokenReuseListRepository {
	return &tokenReuseListRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func (repo *tokenReuseListRepo) Append(ctx context.Context, jit, userID uuid.UUID, expiresAt time.Time) error {
	ctx, span := repo.tracer.Start(ctx, "TokenReuseListRepo.Append")
	defer span.End()

	err := repo.queries(ctx).TokenReuseListAppend(ctx, sqlc2.TokenReuseListAppendParams{
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

	exists, err := repo.queries(ctx).TokenReuseListExists(ctx, sqlc2.TokenReuseListExistsParams{
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
