package security

import (
	"context"
	"time"

	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/shared/ports"
	"lib/database"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type tokenReuseListRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.TokenReuseListRepository = (*tokenReuseListRepo)(nil)

func NewTokenReuseRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) ports.TokenReuseListRepository {
	return &tokenReuseListRepo{
		q:      q,
		log:    l,
		tracer: tracer,
		dbe:    database.NewErrorHandler("token"),
	}
}

func (repo *tokenReuseListRepo) Append(ctx context.Context, jit, userID uuid.UUID, expiresAt time.Time) error {
	ctx, span := repo.tracer.Start(ctx, "Append")
	defer span.End()
	err := database.Queries(ctx, repo.q).TokenReuseListAppend(ctx, sqlc.TokenReuseListAppendParams{
		Jit:       jit,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}

func (repo *tokenReuseListRepo) Exists(ctx context.Context, jit, userID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "Exists")
	defer span.End()
	exists, err := database.Queries(ctx, repo.q).TokenReuseListExists(ctx, sqlc.TokenReuseListExistsParams{
		Jit:    jit,
		UserID: userID,
	})
	if err != nil {
		return false, repo.dbe(err)
	}
	return exists, nil
}

func (repo *tokenReuseListRepo) ClearExpired(ctx context.Context) error {
	ctx, span := repo.tracer.Start(ctx, "ClearExpired")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteExpiredTokenReuseListEntries(ctx)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
