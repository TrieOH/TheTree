package ports

import (
	"context"
	"uuid"

	"Informd/models"
)

type StepRepo interface {
	Create(ctx context.Context, toCreate models.Step) (*models.Step, error)
	List(ctx context.Context, formID uuid.UUID) ([]models.Step, error)
	BulkEdit(ctx context.Context, steps []models.Step) error
}
