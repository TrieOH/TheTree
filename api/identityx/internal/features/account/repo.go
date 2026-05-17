package account

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/shared/ports"
	"context"
	"lib/database"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type accountRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.AccountRepository = (*accountRepo)(nil)

func NewRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) ports.AccountRepository {
	return &accountRepo{
		q:      q,
		log:    l,
		tracer: tracer,
		dbe:    database.NewErrorHandler("user"),
	}
}

func (repo *accountRepo) Verify(ctx context.Context, userID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "Verify")
	span.SetAttributes(attribute.String("user.id", userID.String()))
	defer span.End()
	wasVerified, err := database.Queries(ctx, repo.q).VerifyUser(ctx, userID)
	if err != nil {
		return false, repo.dbe(err)
	}
	span.SetAttributes(attribute.Bool("user.was_already_verified", !wasVerified))
	return !wasVerified, nil
}

func (repo *accountRepo) ResetPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error {
	ctx, span := repo.tracer.Start(ctx, "ResetPassword")
	span.SetAttributes(attribute.String("user.id", userID.String()))
	defer span.End()
	err := database.Queries(ctx, repo.q).ResetUserPassword(ctx, sqlc.ResetUserPasswordParams{
		PasswordHash: string(passwordHash),
		ID:           userID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
