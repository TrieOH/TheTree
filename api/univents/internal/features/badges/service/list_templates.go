package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error) {
	return o.Repo.ListByEdition(ctx, editionID)
}
