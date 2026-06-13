package ports

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

type StepRepo interface {
	Create(ctx context.Context, toCreate models.Step) (*models.Step, error)
	List(ctx context.Context, formID uuid.UUID) ([]models.Step, error)
	BulkEdit(ctx context.Context, steps []models.Step) error
}
