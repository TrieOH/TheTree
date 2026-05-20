package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(ctx context.Context, toCreate models.Project) (*models.Project, error)
	GetByIDExternal(ctx context.Context, projectID, ownerID uuid.UUID) (*models.Project, error)
	GetByIDInternal(ctx context.Context, projectID uuid.UUID) (*models.Project, error)
	IsOwnerOf(ctx context.Context, projectID, ownerID uuid.UUID) (bool, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]models.Project, error)
	Update(ctx context.Context, toUpdate models.Project, ownerID uuid.UUID) (*models.Project, error)
	Delete(ctx context.Context, projectID, ownerID uuid.UUID) error
}
