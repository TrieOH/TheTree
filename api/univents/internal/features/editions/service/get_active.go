package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetActive(ctx context.Context, eventID uuid.UUID) (*models.Edition, error) {
	return o.Editions.GetActive(ctx, eventID)
}
