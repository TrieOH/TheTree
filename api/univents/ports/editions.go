package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type EditionRepo interface {
	Create(ctx context.Context, toCreate *models.Edition) (*models.Edition, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Edition, error)
	GetBySlug(ctx context.Context, eventID uuid.UUID, slug string) (*models.Edition, error)
	ListPublic(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error)
	ListDraft(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error)
	Publish(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, edition *models.Edition) (*models.Edition, error)
}
