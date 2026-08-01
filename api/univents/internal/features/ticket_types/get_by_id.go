package ticket_types

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.TicketType, error) {
	ctx, span := telemetry.StartSpan(ctx, "TicketTypesService.GetByID")
	defer span.End()
	return o.ticketTypes.GetByID(ctx, id)
}
