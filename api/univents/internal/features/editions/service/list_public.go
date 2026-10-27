package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListPublic(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return o.Editions.ListPublic(ctx, event.ID)
}
