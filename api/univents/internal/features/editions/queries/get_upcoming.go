package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetUpcoming(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := q.tracer.Start(ctx, "EditionService.GetUpcoming")
	defer span.End()
	return q.editions.GetUpcoming(ctx, eventID)
}
