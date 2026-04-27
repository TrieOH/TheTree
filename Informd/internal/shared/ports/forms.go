package ports

import (
	"Informd/internal/shared/contracts"
	"context"

	"github.com/google/uuid"
)

type FormsRepo interface {
	Create(ctx context.Context, toCreate contracts.Form) (*contracts.Form, error)
	GetByID(ctx context.Context, id uuid.UUID) (*contracts.Form, error)
	BulkGet(ctx context.Context, ids []uuid.UUID) ([]contracts.Form, error)
}
