package ports

import (
	"IdentityX/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(ctx context.Context, toCreate contracts.Project) (*contracts.Project, error)
	GetByIDExternal(ctx context.Context, projectID, ownerID uuid.UUID) (*contracts.Project, error)
	GetByIDInternal(ctx context.Context, projectID uuid.UUID) (*contracts.Project, error)
	IsOwnerOf(ctx context.Context, projectID, ownerID uuid.UUID) (bool, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]contracts.Project, error)
	Update(ctx context.Context, toUpdate contracts.Project, ownerID uuid.UUID) (*contracts.Project, error)
	Delete(ctx context.Context, projectID, ownerID uuid.UUID) error
}
