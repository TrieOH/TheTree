package editions

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) GetByEventAndEditionSlug(ctx context.Context, eventSlug, editionSlug string) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.GetByEventAndEditionSlug")
	defer span.End()

	event, err := o.events.GetBySlug(ctx, eventSlug)
	if err != nil {
		return nil, err
	}

	return o.editions.GetBySlug(ctx, event.ID, editionSlug)
}
