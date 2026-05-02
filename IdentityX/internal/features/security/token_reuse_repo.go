package security

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"
	"time"

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
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ ports.TokenReuseListRepository = (*tokenReuseListRepo)(nil)

func NewTokenReuseRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) ports.TokenReuseListRepository {
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
		return errx.DB(err, "token")
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
		return false, errx.DB(err, "token")
	}
	return exists, nil
}

func (repo *tokenReuseListRepo) ClearExpired(ctx context.Context) error {
	ctx, span := repo.tracer.Start(ctx, "TokenReuseListRepo.ClearExpired")
	defer span.End()

	err := repo.queries(ctx).DeleteExpiredTokenReuseListEntries(ctx)
	if err != nil {
		return errx.DB(err, "token")
	}
	return nil
}
