package ports

import (
	"context"

	"payssage/models"

	"github.com/google/uuid"
)

type ApiKeysRepo interface {
	Create(ctx context.Context, toCreate models.APIKey) (*models.APIKey, error)
	GetByPrefix(ctx context.Context, prefix string) ([]models.APIKey, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]models.APIKey, error)
	Revoke(ctx context.Context, id, workspaceID uuid.UUID) (*models.APIKey, error)
}
