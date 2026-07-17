package ports

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

type CollectorRepo interface {
	Create(ctx context.Context, toCreate models.Collector) (*models.Collector, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Collector, error)
	List(ctx context.Context) ([]models.Collector, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Collector, error)
	ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Collector, error)
	Revoke(ctx context.Context, id uuid.UUID) error
}
