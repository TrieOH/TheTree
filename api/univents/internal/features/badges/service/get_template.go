package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) GetTemplate(ctx context.Context, templateID uuid.UUID) (*models.BadgeTemplate, error) {
	return o.Repo.GetByID(ctx, templateID)
}
