package ticket_types

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.TicketTypes
}

func New(ops *services.TicketTypes) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
