package editions

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetActive(ctx context.Context, eventID uuid.UUID) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.GetActive")
	defer span.End()
	return o.editions.GetActive(ctx, eventID)
}
