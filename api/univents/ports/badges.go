package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type BadgeTemplateRepo interface {
	CreateTemplate(ctx context.Context, template *models.BadgeTemplate) (*models.BadgeTemplate, error)
	GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.BadgeTemplate, error)
	ListTemplatesByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error)
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
}
