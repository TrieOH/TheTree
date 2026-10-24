package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.TicketType, error) {
	ctx, span := telemetry.StartSpan(ctx, "TicketTypesService.ListByEdition")
	defer span.End()
	return o.TicketTypes.ListByEdition(ctx, editionID)
}
