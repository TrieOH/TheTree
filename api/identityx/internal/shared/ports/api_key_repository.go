package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type ApiKeyRepository interface {
	Upsert(ctx context.Context, key models.ApiKey) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*models.ApiKey, error)
	Delete(ctx context.Context, projectID uuid.UUID) error
}
