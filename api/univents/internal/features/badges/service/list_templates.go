package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error) {
	return o.Repo.ListByEdition(ctx, editionID)
}
