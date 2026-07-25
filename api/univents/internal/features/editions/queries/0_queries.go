package queries

import (
	"univents/ports"
)

type Queries struct {
	events   ports.EventRepo
	editions ports.EditionRepo
}

func NewQueries(
	events ports.EventRepo,
	editions ports.EditionRepo,
) *Queries {
	return &Queries{
		events:   events,
		editions: editions,
	}
}
