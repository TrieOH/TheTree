package ports

import (
	"Informd/models"
	"context"

	"github.com/google/uuid"
)

type ApiKeysRepo interface {
	Create(ctx context.Context, toCreate models.APIKey) (*models.APIKey, error)
	GetByPrefix(ctx context.Context, prefix string) ([]models.APIKey, error)
	Revoke(ctx context.Context, id, userID uuid.UUID) (*models.APIKey, error)
}
