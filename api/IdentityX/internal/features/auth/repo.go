package auth

import (
	"IdentityX/contracts"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/ports"
	"context"
	"lib/database"
	"lib/errx"
	"lib/xslices"

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
	dbe    *errx.DBHandler
}

func (repo *userRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func (repo *userRepo) span(ctx context.Context, op string) (context.Context, trace.Span) {
	return repo.tracer.Start(ctx, "UserRepo."+op)
}

var _ ports.UserRepository = (*userRepo)(nil)

func NewRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer, dbe *errx.DBHandler) ports.UserRepository {
	return &userRepo{
		q:      q,
		log:    l,
		tracer: tracer,
		dbe:    dbe,
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
	ctx, span := repo.span(ctx, "Register")
	defer span.End()
	sqlcUser, err := repo.queries(ctx).RegisterUser(ctx, sqlc.RegisterUserParams{
		Email:        email,
		PasswordHash: password,
		ProjectID:    projectID,
		UserType:     string(userType),
	})
	if err != nil {
		return nil, repo.dbe.DB(err, "user")
	}
	span.SetAttributes(
		attribute.String("user.id", sqlcUser.ID.String()),
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)
	return new(mapUserFromDB(sqlcUser)), nil
}

func (repo *userRepo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	ctx, span := repo.span(ctx, "UpdateLastLogin")
	defer span.End()
	err := repo.queries(ctx).UpdateUserLastLogin(ctx, userID)
	if err != nil {
		return repo.dbe.DB(err, "user")
	}
	return nil
}

func (repo *userRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (*contracts.User, error) {
	ctx, span := repo.span(ctx, "GetUserByID")
	span.SetAttributes(attribute.String("user.id", userID.String()))
	defer span.End()
	sqlcUser, err := repo.queries(ctx).GetUserById(ctx, userID)
	if err != nil {
		return nil, repo.dbe.DB(err, "user")
	}
	span.SetAttributes(
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)
	return new(mapUserFromDB(sqlcUser)), nil
}

func (repo *userRepo) GetUserByEmail(ctx context.Context, email string, projectID *uuid.UUID) (*contracts.User, error) {
	ctx, span := repo.span(ctx, "GetUserByEmail")
	defer span.End()
	sqlcUser, err := repo.queries(ctx).GetUserByEmail(ctx, sqlc.GetUserByEmailParams{
		Email:     email,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, repo.dbe.DB(err, "user")
	}
	span.SetAttributes(
		attribute.String("user.id", sqlcUser.ID.String()),
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)
	return new(mapUserFromDB(sqlcUser)), nil
}

func (repo *userRepo) ListFromProject(ctx context.Context, projectID uuid.UUID) ([]contracts.User, error) {
	ctx, span := repo.span(ctx, "ResetPassword")
	defer span.End()
	sqlcUsers, err := repo.queries(ctx).ListUsersFromProject(ctx, &projectID)
	if err != nil {
		return nil, repo.dbe.DB(err, "user")
	}
	return xslices.MapSlice(sqlcUsers, mapUserFromDB), nil
}

func (repo *userRepo) GetByIDFromProject(ctx context.Context, userID, projectID uuid.UUID) (*contracts.User, error) {
	ctx, span := repo.span(ctx, "GetUserByIDFromProject")
	defer span.End()
	sqlcUser, err := repo.queries(ctx).GetUserByIDFromProject(ctx, sqlc.GetUserByIDFromProjectParams{
		ID:        userID,
		ProjectID: &projectID,
	})
	if err != nil {
		return nil, repo.dbe.DB(err, "user")
	}
	return new(mapUserFromDB(sqlcUser)), nil
}
