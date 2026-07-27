package queries

import (
	"univents/ports"
)

type Queries struct {
	events ports.EventRepo
}

func NewQueries(
	events ports.EventRepo,
) *Queries {
	return &Queries{
		events: events,
	}
}
