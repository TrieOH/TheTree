package service

import (
	"context"
	"univents/models"
)

func (o *Operations) GetByEventAndEditionSlug(ctx context.Context, eventSlug, editionSlug string) (*models.Edition, error) {
	event, err := o.Events.GetBySlug(ctx, eventSlug)
	if err != nil {
		return nil, err
	}

	return o.Editions.GetBySlug(ctx, event.ID, editionSlug)
}
