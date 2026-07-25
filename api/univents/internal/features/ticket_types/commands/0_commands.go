package commands

import (
	"univents/ports"
)

type Commands struct {
	events      ports.EventRepo
	editions    ports.EditionRepo
	ticketTypes ports.TicketTypeRepo
}

func NewCommands(
	events ports.EventRepo,
	editions ports.EditionRepo,
	ticketTypes ports.TicketTypeRepo,
) *Commands {
	return &Commands{
		events:      events,
		editions:    editions,
		ticketTypes: ticketTypes,
	}
}
