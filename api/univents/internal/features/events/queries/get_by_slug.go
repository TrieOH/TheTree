package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (q *Queries) GetBySlug(ctx context.Context, slug string) (*models.Event, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventService.GetBySlug")
	defer span.End()
	return q.events.GetBySlug(ctx, slug)
}
