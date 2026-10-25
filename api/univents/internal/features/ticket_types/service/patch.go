package service

import (
	"context"
	"univents/models"
)

func (o *Operations) Patch(ctx context.Context, payload models.PatchTicketTypeInput) (*models.TicketType, error) {
	existing, err := o.TicketTypes.GetByID(ctx, payload.TicketTypeID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	ticketType := &models.TicketType{
		Name:        payload.Name,
		Description: payload.Description,
		AccessLevel: payload.AccessLevel,
		PriceCents:  payload.PriceCents,
		MaxQuantity: payload.MaxQuantity,
	}

	return o.TicketTypes.Patch(ctx, payload.TicketTypeID, ticketType)
}
