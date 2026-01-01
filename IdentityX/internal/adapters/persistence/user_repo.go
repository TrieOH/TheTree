package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/user"
	"GoAuth/internal/ports/outbound"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ outbound.UserRepository = (*userRepo)(nil)

func NewUserRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.UserRepository {
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
}

func (r userRepo) Register(ctx context.Context, email, password string) (*user.User, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepo.Register")
	defer span.End()

	sqlcUser, err := r.q.RegisterUser(ctx, sqlc.RegisterUserParams{
		Email:        email,
		PasswordHash: password,
	})

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
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

func (r userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepo.GetUserByID",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := r.q.GetUserById(ctx, userID)

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
	}

	span.SetAttributes(
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var usr user.User
	copyUserFromDB(&usr, &sqlcUser)

	return &usr, nil
}

func (r userRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	ctx, span := r.tracer.Start(ctx, "UserRepo.GetUserByEmail")
	defer span.End()

	sqlcUser, err := r.q.GetUserByEmail(ctx, email)

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
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
