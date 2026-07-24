package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetPast(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := q.tracer.Start(ctx, "EditionService.GetPast")
	defer span.End()
	return q.editions.GetPast(ctx, eventID)
}
