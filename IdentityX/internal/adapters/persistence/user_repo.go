package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/domain/user"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *userRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbounds.UserRepository = (*userRepo)(nil)

func NewUserRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbounds.UserRepository {
	return &userRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func copyUserFromDB(dst *user.User, src *sqlc.User) {
	dst.ID = src.ID
	dst.Email = src.Email
	dst.PasswordHash = src.PasswordHash
	dst.UserType = src.UserType
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.IsVerified = src.IsVerified
	dst.VerifiedAt = src.VerifiedAt
}

func (repo *userRepo) Register(ctx context.Context, email, password string) (*user.User, error) {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.Register")
	defer span.End()

	sqlcUser, err := repo.queries(ctx).RegisterUser(ctx, sqlc.RegisterUserParams{
		Email:        email,
		PasswordHash: password,
	})

	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("user.id", sqlcUser.ID.String()),
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var usr user.User
	copyUserFromDB(&usr, &sqlcUser)

	return &usr, nil
}

func (repo *userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.GetUserByID",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := repo.queries(ctx).GetUserById(ctx, userID)

	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var usr user.User
	copyUserFromDB(&usr, &sqlcUser)

	return &usr, nil
}

func (repo *userRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.GetUserByEmail")
	defer span.End()

	sqlcUser, err := repo.queries(ctx).GetUserByEmail(ctx, email)

	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("user.id", sqlcUser.ID.String()),
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var usr user.User
	copyUserFromDB(&usr, &sqlcUser)

	return &usr, nil
}

func (repo *userRepo) Verify(ctx context.Context, userID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.Verify",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	wasVerified, err := repo.queries(ctx).VerifyUser(ctx, userID)
	if err != nil {
		return false, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Bool("user.was_already_verified", !wasVerified))

	return !wasVerified, nil
}
