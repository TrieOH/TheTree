package account

import (
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/ports"
	"context"
	"lib/database"
	"lib/errx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type accountRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
	dbe    *errx.DBHandler
}

func (repo *accountRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ ports.AccountRepository = (*accountRepo)(nil)

func NewRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer, dbe *errx.DBHandler) ports.AccountRepository {
	return &accountRepo{
		q:      q,
		log:    l,
		tracer: tracer,
		dbe:    dbe,
	}
}

func (repo *accountRepo) span(ctx context.Context, op string) (context.Context, trace.Span) {
	return repo.tracer.Start(ctx, "AccountsRepo."+op)
}

func (repo *accountRepo) Verify(ctx context.Context, userID uuid.UUID) (bool, error) {
	ctx, span := repo.span(ctx, "Verify")
	span.SetAttributes(attribute.String("user.id", userID.String()))
	defer span.End()
	wasVerified, err := repo.queries(ctx).VerifyUser(ctx, userID)
	if err != nil {
		return false, repo.dbe.DB(err, "account")
	}
	span.SetAttributes(attribute.Bool("user.was_already_verified", !wasVerified))
	return !wasVerified, nil
}

func (repo *accountRepo) ResetPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error {
	ctx, span := repo.span(ctx, "ResetPassword")
	span.SetAttributes(attribute.String("user.id", userID.String()))
	defer span.End()
	err := repo.queries(ctx).ResetUserPassword(ctx, sqlc.ResetUserPasswordParams{
		PasswordHash: string(passwordHash),
		ID:           userID,
	})
	if err != nil {
		return repo.dbe.DB(err, "account")
	}
	return nil
}
