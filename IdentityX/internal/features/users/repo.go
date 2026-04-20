package users

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
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
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *userRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ ports.UserRepository = (*userRepo)(nil)

func NewRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) ports.UserRepository {
	return &userRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func mapUserFromDB(src *sqlc.User) *contracts.User {
	return &contracts.User{
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
	ctx, span := repo.tracer.Start(ctx, "UserRepo.Register")
	defer span.End()

	sqlcUser, err := repo.queries(ctx).RegisterUser(ctx, sqlc.RegisterUserParams{
		Email:        email,
		PasswordHash: password,
		ProjectID:    projectID,
		UserType:     string(userType),
	})

	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("user.id", sqlcUser.ID.String()),
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	usr := mapUserFromDB(&sqlcUser)
	return usr, nil
}

func (repo *userRepo) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.UpdateLastLogin")
	defer span.End()

	err := repo.queries(ctx).UpdateUserLastLogin(ctx, userID)
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
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

	usr := mapUserFromDB(&sqlcUser)
	return usr, nil
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

	usr := mapUserFromDB(&sqlcUser)
	return usr, nil
}

func (repo *userRepo) GetUserByEmailFromProject(ctx context.Context, email string, projectID uuid.UUID) (*contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.GetUserByEmailFromProject")
	defer span.End()

	sqlcUser, err := repo.queries(ctx).GetUserByEmailFromProject(ctx, sqlc.GetUserByEmailFromProjectParams{
		Email:     email,
		ProjectID: &projectID,
	})

	if err != nil {
		return nil, fail.From(err).WithArgs("user").RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("user.id", sqlcUser.ID.String()),
		attribute.String("user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	usr := mapUserFromDB(&sqlcUser)
	return usr, nil
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

	err := repo.queries(ctx).ResetUserPassword(ctx, sqlc.ResetUserPasswordParams{
		PasswordHash: string(passwordHash),
		ID:           userID,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *userRepo) ListFromProject(ctx context.Context, projectID uuid.UUID) ([]contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.ResetPassword")
	defer span.End()

	sqlcUsers, err := repo.queries(ctx).ListUsersFromProject(ctx, &projectID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	users := make([]contracts.User, 0, len(sqlcUsers))
	for _, sqlcUser := range sqlcUsers {
		user := mapUserFromDB(&sqlcUser)
		users = append(users, *user)
	}

	return users, nil
}

func (repo *userRepo) GetByIDFromProject(ctx context.Context, userID, projectID uuid.UUID) (*contracts.User, error) {
	ctx, span := repo.tracer.Start(ctx, "UserRepo.GetUserByIDFromProject")
	defer span.End()

	sqlcUser, err := repo.queries(ctx).GetUserByIDFromProject(ctx, sqlc.GetUserByIDFromProjectParams{
		ID:        userID,
		ProjectID: &projectID,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	return mapUserFromDB(&sqlcUser), nil
}
