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

type ProjectRepo interface {
	Create(ctx context.Context, project models.Project) (*models.Project, error)

	GetByID(ctx context.Context, projectID, ownerID uuid.UUID) (*models.Project, error)
	GetKeysByID(ctx context.Context, projectID, ownerID uuid.UUID) (*models.ProjectKeys, error)
	GetPublicKeyByID(ctx context.Context, projectID uuid.UUID) (string, error)

	List(ctx context.Context, ownerID uuid.UUID) ([]models.Project, error)

	Update(ctx context.Context, project models.Project, ownerID uuid.UUID) (*models.Project, error)
	Delete(ctx context.Context, projectID, ownerID uuid.UUID) error

	// admin section only
	AdminGetByID(ctx context.Context, projectID uuid.UUID) (*models.Project, error)
}

type projectRepo struct {
	q   *sqlc.Queries
	log *zap.Logger
}

func NewProjectRepo(q *sqlc.Queries, log *zap.Logger) ProjectRepo {
	return &projectRepo{
		q:   q,
		log: log,
	}
}

func (r projectRepo) Create(ctx context.Context, project models.Project) (*models.Project, error) {
	if project.OwnerID == uuid.Nil {
		return nil, errors.New("owner_id is required")
	}
	if project.ProjectName == "" {
		return nil, errors.New("project_name is required")
	}

	sqlcProject, err := r.q.CreateProject(ctx, sqlc.CreateProjectParams{
		ProjectName: project.ProjectName,
		OwnerID:     project.OwnerID,
		Metadata:    project.Metadata,
		IsActive:    project.IsActive,
		PubKey:      project.PubKey,
		PrivKey:     string(project.PrivKey),
	})
	if err != nil {
		return nil, err
	}

	var out models.Project
	if err := copier.Copy(&out, &sqlcProject); err != nil {
		r.log.Error("failed to copy project", zap.Error(err))
		return nil, fmt.Errorf("failed to copy project: %w", err)
	}

	return &out, nil
}

func (r projectRepo) GetByID(ctx context.Context, projectID, ownerID uuid.UUID) (*models.Project, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project_id is required")
	}

	sqlcProject, err := r.q.GetProjectById(ctx, sqlc.GetProjectByIdParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
	if err != nil {
		return nil, err
	}

	var project models.Project
	if err := copier.Copy(&project, &sqlcProject); err != nil {
		return nil, err
	}

	return &project, nil
}

func (r projectRepo) GetKeysByID(ctx context.Context, projectID, ownerID uuid.UUID) (*models.ProjectKeys, error) {
	if projectID == uuid.Nil {
		return nil, errors.New("project_id is required")
	}

	keys, err := r.q.GetProjectKeysById(ctx, sqlc.GetProjectKeysByIdParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
	if err != nil {
		return nil, err
	}

	return &models.ProjectKeys{
		PubKey:  keys.PubKey,
		PrivKey: []byte(keys.PrivKey),
	}, nil
}

func (r projectRepo) GetPublicKeyByID(ctx context.Context, projectID uuid.UUID) (string, error) {
	if projectID == uuid.Nil {
		return "", errors.New("project_id is required")
	}

	pub, err := r.q.GetProjectPublicKeyById(ctx, projectID)
	if err != nil {
		return "", err
	}

	return pub, nil
}

func (r projectRepo) List(ctx context.Context, ownerID uuid.UUID) ([]models.Project, error) {
	if ownerID == uuid.Nil {
		return nil, errors.New("owner_id is required")
	}

	sqlcProjects, err := r.q.ListProjects(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	var projects []models.Project
	for _, p := range sqlcProjects {
		var project models.Project
		if err := copier.Copy(&project, &p); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (r projectRepo) Update(ctx context.Context, project models.Project, ownerID uuid.UUID) (*models.Project, error) {
	if project.ID == uuid.Nil {
		return nil, errors.New("project_id is required")
	}

	sqlcProject, err := r.q.UpdateProject(ctx, sqlc.UpdateProjectParams{
		ID:          project.ID,
		OwnerID:     ownerID,
		ProjectName: project.ProjectName,
		Metadata:    project.Metadata,
	})
	if err != nil {
		return nil, err
	}

	if err := copier.Copy(&project, &sqlcProject); err != nil {
		return nil, err
	}

	return &project, nil
}

func (r projectRepo) Delete(ctx context.Context, projectID, ownerID uuid.UUID) error {
	if projectID == uuid.Nil {
		return errors.New("project_id is required")
	}

	return r.q.DeleteProject(ctx, sqlc.DeleteProjectParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
}

func (r projectRepo) AdminGetByID(ctx context.Context, projectID uuid.UUID) (*models.Project, error) {
	sqlcProject, err := r.q.AdminGetProjectById(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var project models.Project
	if err := copier.Copy(&project, &sqlcProject); err != nil {
		return nil, err
	}

	return &project, nil
}
