package auth

import (
	"IdentityX/contracts"
	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/shared/ports"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type userRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.UserRepository = (*userRepo)(nil)

func NewRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) ports.UserRepository {
	return &userRepo{
		q:      q,
		log:    l,
		tracer: tracer,
		dbe:    database.NewErrorHandler("user"),
	}
}

func mapUserFromDB(src sqlc.User) contracts.User {
	return contracts.User{
		ID:           src.ID,
		UserType:     contracts.UserType(src.UserType),
		ProjectID:    src.ProjectID,
		Email:        src.Email,
		PasswordHash: src.PasswordHash,
		LastLoginAt:  src.LastLoginAt,
		IsVerified:   src.IsVerified,
		VerifiedAt:   src.VerifiedAt,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}
}

func (repo *userRepo) Register(ctx context.Context, email, password string, projectID *uuid.UUID, userType contracts.UserType) (*contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "Register")
	defer span.End()
	sqlcUser, err := database.Queries(ctx, repo.q).RegisterUser(ctx, sqlc.RegisterUserParams{
		Email:        email,
		PasswordHash: password,
		ProjectID:    projectID,
		UserType:     string(userType),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	span.SetAttributes(
		attribute.String("user.id", sqlcUser.ID.String()),
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)
	return new(mapUserFromDB(sqlcUser)), nil
}

func (repo *userRepo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "UpdateLastLogin")
	defer span.End()
	err := database.Queries(ctx, repo.q).UpdateUserLastLogin(ctx, userID)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}

func (repo *userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "GetUserByID")
	span.SetAttributes(attribute.String("user.id", userID.String()))
	defer span.End()
	sqlcUser, err := database.Queries(ctx, repo.q).GetUserById(ctx, userID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	span.SetAttributes(
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)
	return new(mapUserFromDB(sqlcUser)), nil
}

func (repo *userRepo) GetUserByEmail(ctx context.Context, email string, projectID *uuid.UUID) (*contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "GetUserByEmail")
	defer span.End()
	sqlcUser, err := database.Queries(ctx, repo.q).GetUserByEmail(ctx, sqlc.GetUserByEmailParams{
		Email:     email,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	span.SetAttributes(
		attribute.String("user.id", sqlcUser.ID.String()),
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)
	return new(mapUserFromDB(sqlcUser)), nil
}

func (repo *userRepo) ListFromProject(ctx context.Context, projectID uuid.UUID) ([]contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "ResetPassword")
	defer span.End()
	sqlcUsers, err := database.Queries(ctx, repo.q).ListUsersFromProject(ctx, &projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcUsers, mapUserFromDB), nil
}

func (repo *userRepo) GetByIDFromProject(ctx context.Context, userID, projectID uuid.UUID) (*contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "GetUserByIDFromProject")
	defer span.End()
	sqlcUser, err := database.Queries(ctx, repo.q).GetUserByIDFromProject(ctx, sqlc.GetUserByIDFromProjectParams{
		ID:        userID,
		ProjectID: &projectID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapUserFromDB(sqlcUser)), nil
}
