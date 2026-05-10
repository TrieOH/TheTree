package ports

import (
	"context"
	"payssage/contracts"

	"github.com/google/uuid"
)

type ApiKeysRepo interface {
	Create(ctx context.Context, toCreate contracts.APIKey) (*contracts.APIKey, error)
	GetByPrefix(ctx context.Context, prefix string) ([]contracts.APIKey, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]contracts.APIKey, error)
	Revoke(ctx context.Context, id, workspaceID uuid.UUID) (*contracts.APIKey, error)
}
