package ports

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

type ResponseRepo interface {
	Create(ctx context.Context, toCreate models.Response) (*models.Response, error)
	Finish(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Response, error)
	ListByForm(ctx context.Context, formID uuid.UUID) ([]models.Response, error)
}
