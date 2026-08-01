package ticket_types

import (
	"univents/ports"
)

type Operations struct {
	events      ports.EventRepo
	editions    ports.EditionRepo
	ticketTypes ports.TicketTypeRepo
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	ticketTypes ports.TicketTypeRepo,
) *Operations {
	return &Operations{
		events:      events,
		editions:    editions,
		ticketTypes: ticketTypes,
	}
}
