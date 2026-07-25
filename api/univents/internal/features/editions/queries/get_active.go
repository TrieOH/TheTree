package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetActive(ctx context.Context, eventID uuid.UUID) (*models.Edition, error) {
	ctx, span := telemetry.StartSpan(ctx, "EditionService.GetActive")
	defer span.End()
	return q.editions.GetActive(ctx, eventID)
}
