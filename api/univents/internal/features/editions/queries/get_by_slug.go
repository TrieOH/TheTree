package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (q *Queries) GetByEventAndEditionSlug(ctx context.Context, eventSlug, editionSlug string) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.GetByEventAndEditionSlug")
	defer span.End()

	event, err := q.events.GetBySlug(ctx, eventSlug)
	if err != nil {
		return nil, err
	}

	return q.editions.GetBySlug(ctx, event.ID, editionSlug)
}
