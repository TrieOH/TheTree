package products

import (
	"univents/ports"
)

type Operations struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	products ports.ProductRepo
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	products ports.ProductRepo,
) *Operations {
	return &Operations{
		events:   events,
		editions: editions,
		products: products,
	}
}
