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

type ProjectUserRepo interface {
	Register(ctx context.Context, user models.ProjectUser) (*models.ProjectUser, error)

	GetByIDExternal(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) (*models.ProjectUser, error)
	GetByIDInternal(ctx context.Context, projectUserID, projectID uuid.UUID) (*models.ProjectUser, error)

	GetByEmailExternal(ctx context.Context, projectID uuid.UUID, email string, ownerID uuid.UUID) (*models.ProjectUser, error)
	GetByEmailInternal(ctx context.Context, projectID uuid.UUID, email string) (*models.ProjectUser, error)

	ListExternal(ctx context.Context, projectID, ownerID uuid.UUID) ([]models.ProjectUser, error)
	ListInternal(ctx context.Context, projectID uuid.UUID) ([]models.ProjectUser, error)

	Update(ctx context.Context, user models.ProjectUser, ownerID uuid.UUID) (*models.ProjectUser, error)
	Delete(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) error
}

type projectUserRepo struct {
	q   *sqlc.Queries
	log *zap.Logger
}

func NewProjectUserRepo(q *sqlc.Queries, log *zap.Logger) ProjectUserRepo {
	return &projectUserRepo{
		q:   q,
		log: log,
	}
}

func mapProjectUserFromDB(dst *models.ProjectUser, src *sqlc.ProjectUser) {
	dst.ID = src.ID
	dst.ProjectID = src.ProjectID
	dst.Email = src.Email
	dst.UserType = src.UserType
	dst.Metadata = &src.Metadata
	dst.IsActive = src.IsActive
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.LastLoginAt = src.LastLoginAt
}

func (r projectUserRepo) Register(ctx context.Context, user models.ProjectUser) (*models.ProjectUser, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectUserRepo.Register",
		trace.WithAttributes(
			attribute.String("user.project_id", user.ProjectID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := r.q.RegisterProjectUser(ctx, sqlc.RegisterProjectUserParams{
		ProjectID:    user.ProjectID,
		Email:        user.Email,
		PasswordHash: user.Password,
		Metadata:     *user.Metadata,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.String("project_user.id", sqlcUser.ID.String()),
		attribute.String("project_user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	mapProjectUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (r projectUserRepo) GetByIDExternal(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) (*models.ProjectUser, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectUserRepo.GetByIDExternal",
		trace.WithAttributes(
			attribute.String("project_user.id", projectUserID.String()),
			attribute.String("project_user.project_id", projectID.String()),
			attribute.String("project.owner_id", ownerID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := r.q.GetProjectUserById(ctx, sqlc.GetProjectUserByIdParams{
		ID:        projectUserID,
		ProjectID: projectID,
		OwnerID:   ownerID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.String("project_user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
		attribute.Bool("user.is_active", sqlcUser.IsActive),
	)

	var user models.ProjectUser
	mapProjectUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (r projectUserRepo) GetByIDInternal(ctx context.Context, projectUserID, projectID uuid.UUID) (*models.ProjectUser, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectUserRepo.GetByID",
		trace.WithAttributes(
			attribute.String("project_user.project_id", projectID.String()),
			attribute.String("project_user.id", projectUserID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := r.q.GetProjectUserByIdInternal(ctx, sqlc.GetProjectUserByIdInternalParams{
		ID:        projectUserID,
		ProjectID: projectID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
		attribute.Bool("user.is_active", sqlcUser.IsActive),
	)

	var user models.ProjectUser
	mapProjectUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (r projectUserRepo) GetByEmailExternal(ctx context.Context, projectID uuid.UUID, email string, ownerID uuid.UUID) (*models.ProjectUser, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectUserRepo.GetByEmailExternal",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project_user.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := r.q.GetProjectUserByEmailExternal(ctx, sqlc.GetProjectUserByEmailExternalParams{
		ProjectID: projectID,
		Email:     email,
		OwnerID:   ownerID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.String("project_user.id", sqlcUser.ID.String()),
		attribute.String("project_user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var user models.ProjectUser
	mapProjectUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (r projectUserRepo) GetByEmailInternal(ctx context.Context, projectID uuid.UUID, email string) (*models.ProjectUser, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectUserRepo.GetByEmailInternal",
		trace.WithAttributes(
			attribute.String("project_user.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := r.q.GetProjectUserByEmailInternal(ctx, sqlc.GetProjectUserByEmailInternalParams{
		ProjectID: projectID,
		Email:     email,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.String("project_user.id", sqlcUser.ID.String()),
		attribute.String("project_user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	var user models.ProjectUser
	mapProjectUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (r projectUserRepo) ListExternal(ctx context.Context, projectID, ownerID uuid.UUID) ([]models.ProjectUser, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectUserRepo.ListExternal",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcUsers, err := r.q.ListProjectUsersExternal(ctx, sqlc.ListProjectUsersExternalParams{
		ProjectID: projectID,
		OwnerID:   ownerID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.Int("project_users.count", len(sqlcUsers)))

	users := make([]models.ProjectUser, 0, len(sqlcUsers))
	for _, u := range sqlcUsers {
		var user models.ProjectUser
		mapProjectUserFromDB(&user, &u)
		users = append(users, user)
	}

	return users, nil
}

func (r projectUserRepo) ListInternal(ctx context.Context, projectID uuid.UUID) ([]models.ProjectUser, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectUserRepo.ListInternal",
		trace.WithAttributes(
			attribute.String("project.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcUsers, err := r.q.ListProjectUsersInternal(ctx, projectID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.Int("project_users.count", len(sqlcUsers)))

	users := make([]models.ProjectUser, 0, len(sqlcUsers))
	for _, u := range sqlcUsers {
		var user models.ProjectUser
		mapProjectUserFromDB(&user, &u)
		users = append(users, user)
	}

	return users, nil
}

func (r projectUserRepo) Update(ctx context.Context, user models.ProjectUser, ownerID uuid.UUID) (*models.ProjectUser, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectUserRepo.Update",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.project_id", user.ProjectID.String()),
			attribute.String("project_user.id", user.ID.String()),
		),
	)
	defer span.End()

	sqlcUser, err := r.q.UpdateProjectUser(ctx, sqlc.UpdateProjectUserParams{
		ID:           user.ID,
		ProjectID:    user.ProjectID,
		Email:        user.Email,
		PasswordHash: user.Password,
		OwnerID:      ownerID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.String("project_user.type", sqlcUser.UserType),
		attribute.Int64("user.created_at", sqlcUser.CreatedAt.Unix()),
	)

	mapProjectUserFromDB(&user, &sqlcUser)

	return &user, nil
}

func (r projectUserRepo) Delete(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) error {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectUserRepo.Delete",
		trace.WithAttributes(
			attribute.String("project.project_id", projectID.String()),
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project_user.id", projectUserID.String()),
		),
	)
	defer span.End()

	if err := r.q.DeleteProjectUser(ctx, sqlc.DeleteProjectUserParams{
		ID:        projectUserID,
		ProjectID: projectID,
		OwnerID:   ownerID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	span.SetAttributes(attribute.Bool("project_user.deleted", true))

	return nil
}
