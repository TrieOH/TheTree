package ports

import (
	"IdentityX/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type ApiKeyRepository interface {
	Upsert(ctx context.Context, key contracts.ApiKey) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*contracts.ApiKey, error)
	Delete(ctx context.Context, projectID uuid.UUID) error
}
