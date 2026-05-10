package ports

import (
	"Informd/contracts"
	"context"

	"github.com/google/uuid"
)

type ApiKeysRepo interface {
	Create(ctx context.Context, toCreate contracts.APIKey) (*contracts.APIKey, error)
	GetByPrefix(ctx context.Context, prefix string) ([]contracts.APIKey, error)
	BulkGet(ctx context.Context, ids []uuid.UUID) ([]contracts.APIKey, error)
	Revoke(ctx context.Context, id, userID uuid.UUID) (*contracts.APIKey, error)
}
