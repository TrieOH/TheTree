package repo

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/models"
	"GoAuth/internal/sqlc"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type UserRepo interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers() ([]*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID string) error
}

type userRepo struct {
	q   *sqlc.Queries
	log *zap.Logger
}

func NewUserRepo(q *sqlc.Queries, l *zap.Logger) UserRepo {
	return &userRepo{
		q:   q,
		log: l,
	}
}

func copyUserFromDB(dst *models.User, src *sqlc.User) {
	dst.ID = src.ID
	dst.Email = src.Email
	dst.PasswordHash = src.PasswordHash
	dst.UserType = src.UserType
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func (u userRepo) Register(ctx context.Context, email, password string) (*models.User, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "UserRepo.Register")
	defer span.End()

	sqlcUser, err := u.q.RegisterUser(ctx, sqlc.RegisterUserParams{
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

	var user models.User
	copyUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (u userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "UserRepo.GetUserByID",
		trace.WithAttributes(
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := u.q.GetUserById(ctx, userID)

	if err != nil {
		sqlErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlErr)
		return nil, sqlErr
	}

	span.SetAttributes(
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var user models.User
	copyUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (u userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "UserRepo.GetUserByEmail")
	defer span.End()

	sqlcUser, err := u.q.GetUserByEmail(ctx, email)

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

	var user models.User
	copyUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (u userRepo) ListUsers() ([]*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) UpdateUser(_ *models.User) error {
	//TODO implement me
	panic("implement me")
}

func (u userRepo) DeleteUser(_ string) error {
	//TODO implement me
	panic("implement me")
}
