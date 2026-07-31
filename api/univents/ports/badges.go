package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type BadgeTemplateRepo interface {
	Create(ctx context.Context, template *models.BadgeTemplate) (*models.BadgeTemplate, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.BadgeTemplate, error)
	ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
