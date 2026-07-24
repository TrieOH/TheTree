package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListPublic(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	ctx, span := q.tracer.Start(ctx, "EditionService.ListPublic")
	defer span.End()
	event, err := q.events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return q.editions.ListPublic(ctx, event.ID)
}
