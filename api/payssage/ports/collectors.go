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
}
