package ports

import (
	"Informd/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type FormsRepo interface {
	Create(ctx context.Context, toCreate contracts.Form) (*contracts.Form, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]contracts.Form, error)
	ListByNamespace(ctx context.Context, namespaceID *uuid.UUID) ([]contracts.Form, error)
}
