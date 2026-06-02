package ports

import (
	"Informd/models"
	"context"

	"github.com/google/uuid"
)

type FieldsRepo interface {
	Create(ctx context.Context, toCreate models.Field) (*models.Field, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Field, error)
	ListByStepID(ctx context.Context, stepID uuid.UUID) ([]models.Field, error)
	ListByFormID(ctx context.Context, formID uuid.UUID) ([]models.Field, error)
	BulkEdit(ctx context.Context, fields []models.Field) error
	Delete(ctx context.Context, id uuid.UUID) error

	CreateSelectConfig(ctx context.Context, toCreate models.FieldSelectConfig) (*models.FieldSelectConfig, error)
	GetSelectConfig(ctx context.Context, fieldID uuid.UUID) (*models.FieldSelectConfig, error)
	UpdateSelectConfig(ctx context.Context, toUpdate models.FieldSelectConfig) (*models.FieldSelectConfig, error)
	DeleteSelectConfig(ctx context.Context, fieldID uuid.UUID) error
}
