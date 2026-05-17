package ports

import (
	"Informd/models"
	"context"

	"github.com/google/uuid"
)

type NamespaceRepo interface {
	Create(ctx context.Context, toCreate models.Namespace) (*models.Namespace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Namespace, error)
	GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*models.Namespace, error)
	BulkGet(ctx context.Context, ids []uuid.UUID) ([]models.Namespace, error)
}
