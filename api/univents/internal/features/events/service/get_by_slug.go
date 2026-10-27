package service

import (
	"context"
	"univents/models"
)

func (o *Operations) GetBySlug(ctx context.Context, slug string) (*models.Event, error) {
	return o.Events.GetBySlug(ctx, slug)
}
