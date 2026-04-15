package users

import (
	"IdentityX/internal/platform/database"
	sqlc2 "IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/ports"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userRepo struct {
	q      *sqlc2.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *userRepo) queries(ctx context.Context) *sqlc2.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ ports.UserRepository = (*userRepo)(nil)

func NewRepo(q *sqlc2.Queries, l *zap.Logger, tracer trace.Tracer) ports.UserRepository {
	return &userRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func copyUserFromDB(dst *contracts.User, src *sqlc2.User) {
	dst.ID = src.ID
	dst.Email = src.Email
	dst.PasswordHash = src.PasswordHash
	dst.UserType = src.UserType
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.IsVerified = src.IsVerified
	dst.VerifiedAt = src.VerifiedAt
}

func (repo *userRepo) Register(ctx context.Context, email, password string) (*contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.Register")
	defer span.End()

	sqlcUser, err := repo.queries(ctx).RegisterUser(ctx, sqlc2.RegisterUserParams{
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

	var usr contracts.User
	copyUserFromDB(&usr, &sqlcUser)

	return &usr, nil
}

func (repo *userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.GetUserByID",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := repo.queries(ctx).GetUserById(ctx, userID)

	if err != nil {
		return nil, fail.From(err).WithArgs("user").RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var usr contracts.User
	copyUserFromDB(&usr, &sqlcUser)

	return &usr, nil
}

func (repo *userRepo) GetUserByEmail(ctx context.Context, email string) (*contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.GetUserByEmail")
	defer span.End()

	sqlcUser, err := repo.queries(ctx).GetUserByEmail(ctx, email)

	if err != nil {
		return nil, fail.From(err).WithArgs("user").RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("user.id", sqlcUser.ID.String()),
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var usr contracts.User
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

func (repo *userRepo) ResetPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.ResetPassword",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	err := repo.queries(ctx).ResetUserPassword(ctx, sqlc2.ResetUserPasswordParams{
		PasswordHash: string(passwordHash),
		ID:           userID,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}
