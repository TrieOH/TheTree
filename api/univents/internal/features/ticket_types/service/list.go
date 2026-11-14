package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.TicketType, error) {
	return o.TicketTypes.ListByEdition(ctx, editionID)
}
