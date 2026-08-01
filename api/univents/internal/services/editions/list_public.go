package editions

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListPublic(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.ListPublic")
	defer span.End()
	event, err := o.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return o.editions.ListPublic(ctx, event.ID)
}
