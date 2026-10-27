package service

import (
	"context"
	"univents/models"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateTicketTypeInput) (*models.TicketType, error) {
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
