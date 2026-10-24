package service

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateTicketTypeInput) (*models.TicketType, error) {
	ctx, span := telemetry.StartSpan(ctx, "TicketTypesService.Create")
	defer span.End()

	edition, err := o.Editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	ticketType := &models.TicketType{
		EditionID:   payload.EditionID,
		Name:        payload.Name,
		Description: payload.Description,
		AccessLevel: payload.AccessLevel,
		PriceCents:  payload.PriceCents,
		MaxQuantity: payload.MaxQuantity,
	}

	return o.TicketTypes.Create(ctx, ticketType)
}
