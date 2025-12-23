package repo

import (
	"GoAuth/internal/models"
	"GoAuth/internal/sqlc"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
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

func copyProjectUserFromDB(dst *models.ProjectUser, src *sqlc.ProjectUser) error {
	return copier.Copy(dst, src)
}

func (r projectUserRepo) Register(ctx context.Context, user models.ProjectUser) (*models.ProjectUser, error) {
	if user.ProjectID == uuid.Nil {
		return nil, errors.New("project_id is required")
	}
	if user.Email == "" {
		return nil, errors.New("email is required")
	}

	sqlcUser, err := r.q.RegisterProjectUser(ctx, sqlc.RegisterProjectUserParams{
		ProjectID:    user.ProjectID,
		Email:        user.Email,
		PasswordHash: user.Password,
		Metadata:     *user.Metadata,
	})
	if err != nil {
		return nil, err
	}

	if err := copyProjectUserFromDB(&user, &sqlcUser); err != nil {
		r.log.Error("failed to copy project user", zap.Error(err))
		return nil, fmt.Errorf("failed to copy project user: %w", err)
	}

	return &user, nil
}

func (r projectUserRepo) GetByIDExternal(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) (*models.ProjectUser, error) {
	sqlcUser, err := r.q.GetProjectUserById(ctx, sqlc.GetProjectUserByIdParams{
		ID:        projectUserID,
		ProjectID: projectID,
		OwnerID:   ownerID,
	})
	if err != nil {
		return nil, err
	}

	var user models.ProjectUser
	if err := copyProjectUserFromDB(&user, &sqlcUser); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r projectUserRepo) GetByIDInternal(ctx context.Context, projectUserID, projectID uuid.UUID) (*models.ProjectUser, error) {
	sqlcUser, err := r.q.GetProjectUserByIdInternal(ctx, sqlc.GetProjectUserByIdInternalParams{
		ID:        projectUserID,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, err
	}

	var user models.ProjectUser
	if err := copyProjectUserFromDB(&user, &sqlcUser); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r projectUserRepo) GetByEmailExternal(ctx context.Context, projectID uuid.UUID, email string, ownerID uuid.UUID) (*models.ProjectUser, error) {
	sqlcUser, err := r.q.GetProjectUserByEmailExternal(ctx, sqlc.GetProjectUserByEmailExternalParams{
		ProjectID: projectID,
		Email:     email,
		OwnerID:   ownerID,
	})
	if err != nil {
		return nil, err
	}

	var user models.ProjectUser
	if err := copyProjectUserFromDB(&user, &sqlcUser); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r projectUserRepo) GetByEmailInternal(ctx context.Context, projectID uuid.UUID, email string) (*models.ProjectUser, error) {
	sqlcUser, err := r.q.GetProjectUserByEmailInternal(ctx, sqlc.GetProjectUserByEmailInternalParams{
		ProjectID: projectID,
		Email:     email,
	})
	if err != nil {
		return nil, err
	}

	var user models.ProjectUser
	if err := copyProjectUserFromDB(&user, &sqlcUser); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r projectUserRepo) ListExternal(ctx context.Context, projectID, ownerID uuid.UUID) ([]models.ProjectUser, error) {
	sqlcUsers, err := r.q.ListProjectUsersExternal(ctx, sqlc.ListProjectUsersExternalParams{
		ProjectID: projectID,
		OwnerID:   ownerID,
	})
	if err != nil {
		return nil, err
	}

	var users []models.ProjectUser
	for _, u := range sqlcUsers {
		var user models.ProjectUser
		if err := copyProjectUserFromDB(&user, &u); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r projectUserRepo) ListInternal(ctx context.Context, projectID uuid.UUID) ([]models.ProjectUser, error) {
	sqlcUsers, err := r.q.ListProjectUsersInternal(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var users []models.ProjectUser
	for _, u := range sqlcUsers {
		var user models.ProjectUser
		if err := copyProjectUserFromDB(&user, &u); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r projectUserRepo) Update(ctx context.Context, user models.ProjectUser, ownerID uuid.UUID) (*models.ProjectUser, error) {
	if user.ID == uuid.Nil {
		return nil, errors.New("project_user_id is required")
	}

	sqlcUser, err := r.q.UpdateProjectUser(ctx, sqlc.UpdateProjectUserParams{
		ID:           user.ID,
		ProjectID:    user.ProjectID,
		Email:        user.Email,
		PasswordHash: user.Password,
		OwnerID:      ownerID,
	})
	if err != nil {
		return nil, err
	}

	if err := copyProjectUserFromDB(&user, &sqlcUser); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r projectUserRepo) Delete(ctx context.Context, projectUserID, projectID, ownerID uuid.UUID) error {
	return r.q.DeleteProjectUser(ctx, sqlc.DeleteProjectUserParams{
		ID:        projectUserID,
		ProjectID: projectID,
		OwnerID:   ownerID,
	})
}
