package ports

import (
	"Informd/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type NamespaceRepo interface {
	Create(ctx context.Context, toCreate contracts.Namespace) (*contracts.Namespace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Namespace, error)
	GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*contracts.Namespace, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]contracts.Namespace, error)
	ListByIDs(ctx context.Context, ids []string) ([]contracts.Namespace, error)
}
