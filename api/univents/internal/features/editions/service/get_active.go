package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) GetActive(ctx context.Context, eventID uuid.UUID) (*models.Edition, error) {
	return o.Editions.GetActive(ctx, eventID)
}
