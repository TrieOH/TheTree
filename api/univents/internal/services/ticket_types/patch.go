package ticket_types

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"
)

func (o *Operations) Patch(ctx context.Context, payload models.PatchTicketTypeInput) (*models.TicketType, error) {
	ctx, span := telemetry.StartSpan(ctx, "TicketTypesService.Patch")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := o.ticketTypes.GetByID(ctx, payload.TicketTypeID)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
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

	return o.ticketTypes.Patch(ctx, payload.TicketTypeID, ticketType)
}
