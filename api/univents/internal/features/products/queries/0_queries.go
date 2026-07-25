package queries

import (
	"univents/ports"
)

type Queries struct {
	editions ports.EditionRepo
	products ports.ProductRepo
}

func NewQueries(
	editions ports.EditionRepo,
	products ports.ProductRepo,
) *Queries {
	return &Queries{
		editions: editions,
		products: products,
	}
}
