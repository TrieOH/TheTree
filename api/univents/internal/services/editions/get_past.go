package editions

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetPast(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.GetPast")
	defer span.End()
	return o.editions.GetPast(ctx, eventID)
}
