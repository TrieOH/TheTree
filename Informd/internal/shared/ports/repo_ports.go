package ports

import (
	"TrieForms/internal/shared/types"
	"context"

	"github.com/google/uuid"
)

type ApiKeysRepo interface {
	Create(ctx context.Context, toCreate types.APIKey) (*types.APIKey, error)
	GetByPrefix(ctx context.Context, prefix string) ([]types.APIKey, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]types.APIKey, error)
	Revoke(ctx context.Context, id, userID uuid.UUID) (*types.APIKey, error)
}

type ProjectsRepo interface {
	Create(ctx context.Context, toCreate types.Project) (*types.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Project, error)
	GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*types.Project, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]types.Project, error)
	ListByIDs(ctx context.Context, ids []string) ([]types.Project, error)
}

type FormsRepo interface {
	Create(ctx context.Context, toCreate types.Form) (*types.Form, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]types.Form, error)
}
