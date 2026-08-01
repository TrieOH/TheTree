package editions

import (
	"univents/ports"
)

type Operations struct {
	events   ports.EventRepo
	editions ports.EditionRepo
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
) *Operations {
	return &Operations{
		events:   events,
		editions: editions,
	}
}
