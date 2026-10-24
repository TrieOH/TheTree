package service

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) GetBySlug(ctx context.Context, slug string) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventService.GetBySlug")
	defer span.End()
	return o.Events.GetBySlug(ctx, slug)
}
