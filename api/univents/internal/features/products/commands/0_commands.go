package commands

import (
	"univents/ports"
)

type Commands struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	products ports.ProductRepo
}

func NewCommands(
	events ports.EventRepo,
	editions ports.EditionRepo,
	products ports.ProductRepo,
) *Commands {
	return &Commands{
		events:   events,
		editions: editions,
		products: products,
	}
}
