package commands

import (
	"univents/ports"
)

type Commands struct {
	events   ports.EventRepo
	editions ports.EditionRepo
}

func NewCommands(
	events ports.EventRepo,
	editions ports.EditionRepo,
) *Commands {
	return &Commands{
		events:   events,
		editions: editions,
	}
}
