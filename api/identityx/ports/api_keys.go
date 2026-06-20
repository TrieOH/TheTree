package ports

import (
	"IdentityX/models"
	"context"
)

type ApiKeysRepo interface {
	Create(ctx context.Context, toCreate models.ApiKey) (*models.ApiKey, error)
	GetByPrefix(ctx context.Context, prefix string) (*models.ApiKey, error)
}
