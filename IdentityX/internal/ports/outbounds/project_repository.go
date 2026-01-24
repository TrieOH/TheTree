package outbounds

import (
	"GoAuth/internal/domain/project"
	"context"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(ctx context.Context, toCreate project.Project) (*project.Project, error)
	GetByIDExternal(ctx context.Context, projectID, ownerID uuid.UUID) (*project.Project, error)
	GetByIDInternal(ctx context.Context, projectID uuid.UUID) (*project.Project, error)
	IsOwnerOf(ctx context.Context, projectID, ownerID uuid.UUID) (bool, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]project.Project, error)
	Update(ctx context.Context, toUpdate project.Project, ownerID uuid.UUID) (*project.Project, error)
	Delete(ctx context.Context, projectID, ownerID uuid.UUID) error
}
