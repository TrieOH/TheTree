package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) GetUpcoming(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	return o.Editions.GetUpcoming(ctx, eventID)
}
