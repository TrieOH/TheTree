package queries

import (
	"univents/ports"
)

type Queries struct {
	editions    ports.EditionRepo
	ticketTypes ports.TicketTypeRepo
}

func NewQueries(
	editions ports.EditionRepo,
	ticketTypes ports.TicketTypeRepo,
) *Queries {
	return &Queries{
		editions:    editions,
		ticketTypes: ticketTypes,
	}
}
