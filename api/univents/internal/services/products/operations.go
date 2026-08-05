package products

import (
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	products ports.ProductRepo
	authz    *authz.Service
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	products ports.ProductRepo,
	authz *authz.Service,
) *Operations {
	return &Operations{
		events:   events,
		editions: editions,
		products: products,
		authz:    authz,
	}
}
