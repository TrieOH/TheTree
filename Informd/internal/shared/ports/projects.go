package ports

import (
	"Informd/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type ProjectsRepo interface {
	Create(ctx context.Context, toCreate contracts.Project) (*contracts.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Project, error)
	GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*contracts.Project, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]contracts.Project, error)
	ListByIDs(ctx context.Context, ids []string) ([]contracts.Project, error)
}
