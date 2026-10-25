package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetPast(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	return o.Editions.GetPast(ctx, eventID)
}
