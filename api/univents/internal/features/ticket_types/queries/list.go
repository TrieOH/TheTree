package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.TicketType, error) {
	ctx, span := telemetry.StartSpan(ctx, "TicketTypesService.ListByEdition")
	defer span.End()
	return q.ticketTypes.ListByEdition(ctx, editionID)
}
