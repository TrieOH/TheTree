package ports

import (
	"Informd/contracts"
	"context"

	"github.com/google/uuid"
)

type NamespaceRepo interface {
	Create(ctx context.Context, toCreate contracts.Namespace) (*contracts.Namespace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Namespace, error)
	GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*contracts.Namespace, error)
	BulkGet(ctx context.Context, ids []uuid.UUID) ([]contracts.Namespace, error)
}
