package ports

import (
	"Informd/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type ApiKeysRepo interface {
	Create(ctx context.Context, toCreate contracts.APIKey) (*contracts.APIKey, error)
	GetByPrefix(ctx context.Context, prefix string) ([]contracts.APIKey, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]contracts.APIKey, error)
	Revoke(ctx context.Context, id, userID uuid.UUID) (*contracts.APIKey, error)
}
