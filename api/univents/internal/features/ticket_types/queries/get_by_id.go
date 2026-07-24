package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*models.TicketType, error) {
	ctx, span := q.tracer.Start(ctx, "TicketTypesService.GetByID")
	defer span.End()
	return q.ticketTypes.GetByID(ctx, id)
}
