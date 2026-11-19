// Package ticket_types is the ticket types feature: its repo, handlers and
// operations. Build wires it.
package ticket_types

import (
	"context"
	"univents/internal/authz"
	"univents/models"
	"univents/ports"
	"uuid"
)

type Deps struct {
	Events      ports.EventRepo
	Editions    ports.EditionRepo
	TicketTypes ports.TicketTypeRepo
	Authz       *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

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

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.TicketType, error) {
	return o.TicketTypes.GetByID(ctx, id)
}

func (o *Operations) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.TicketType, error) {
	return o.TicketTypes.ListByEdition(ctx, editionID)
}

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
