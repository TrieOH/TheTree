package products

import (
	"lib/database"
	"univents/internal/authz"
	"univents/ports"
)

type Operations struct {
	events   ports.EventRepo
	editions ports.EditionRepo
	products ports.ProductRepo
	authz    *authz.Service
	// tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	tx database.TxRunner
}

func NewOperations(
	events ports.EventRepo,
	editions ports.EditionRepo,
	products ports.ProductRepo,
	authz *authz.Service,
	tx database.TxRunner,
) *Operations {
	return &Operations{
		events:   events,
		editions: editions,
		products: products,
		authz:    authz,
		tx:       tx,
	}
}
