package queries

import (
	"context"
	"univents/models"
)

func (q *Queries) GetByEventAndEditionSlug(ctx context.Context, eventSlug, editionSlug string) (*models.Edition, error) {
	ctx, span := q.tracer.Start(ctx, "EditionService.GetByEventAndEditionSlug")
	defer span.End()

	event, err := q.events.GetBySlug(ctx, eventSlug)
	if err != nil {
		return nil, err
	}

	return q.editions.GetBySlug(ctx, event.ID, editionSlug)
}
