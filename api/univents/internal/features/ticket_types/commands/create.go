package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) Create(ctx context.Context, payload models.CreateTicketTypeInput) (*models.TicketType, error) {
	ctx, span := telemetry.StartSpan(ctx, "TicketTypesService.Create")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	edition, err := c.editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	member, err := c.events.GetMember(ctx, edition.EventID, ident.Sub.ID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}
	if err != nil {
		return nil, err
	}
	if !member.Role.Minimum(models.EventMemberRoleAdmin) {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	ticketType := &models.TicketType{
		EditionID:   payload.EditionID,
		Name:        payload.Name,
		Description: payload.Description,
		AccessLevel: payload.AccessLevel,
		PriceCents:  payload.PriceCents,
		MaxQuantity: payload.MaxQuantity,
	}

	return c.ticketTypes.Create(ctx, ticketType)
}
