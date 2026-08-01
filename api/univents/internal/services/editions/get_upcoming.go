package editions

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetUpcoming(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.GetUpcoming")
	defer span.End()
	return o.editions.GetUpcoming(ctx, eventID)
}
