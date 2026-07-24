package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type EditionRepo interface {
	Create(ctx context.Context, toCreate *models.Edition) (*models.Edition, error)
	ListPublic(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error)
	ListDraft(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error)
}
