package ports

import (
	"Informd/models"
	"context"

	"github.com/google/uuid"
)

type FormsRepo interface {
	Create(ctx context.Context, toCreate models.Form) (*models.Form, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Form, error)
	BulkGet(ctx context.Context, ids []uuid.UUID, params models.BulkGetParams) ([]models.Form, error)
}
