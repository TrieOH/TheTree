package ports

import (
	"IdentityX/models"
	"context"
)

type ApiKeysRepo interface {
	Create(ctx context.Context, toCreate models.APIKey) (*models.APIKey, error)
	GetByPrefix(ctx context.Context, prefix string) (*models.APIKey, error)
}
