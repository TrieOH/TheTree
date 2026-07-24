package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.TicketType, error) {
	ctx, span := q.tracer.Start(ctx, "TicketTypesService.ListByEdition")
	defer span.End()
	return q.ticketTypes.ListByEdition(ctx, editionID)
}
