package ticket_types

import (
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	events      ports.EventRepo
	editions    ports.EditionRepo
	ticketTypes ports.TicketTypeRepo
	authz       *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	ticketTypes ports.TicketTypeRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events:      events,
		editions:    editions,
		ticketTypes: ticketTypes,
		authz:       authz,
	}
}
