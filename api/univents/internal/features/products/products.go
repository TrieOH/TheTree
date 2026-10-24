// Package products is the products feature: the layers live in the
// subpackages, and Build wires the feature from here.
package products

import (
	"univents/internal/features/products/handlers"
	"univents/internal/features/products/repo"
	"univents/internal/features/products/service"
	"univents/internal/sqlc"
)

type (
	Repo       = repo.Repo
	Deps       = service.Deps
	Operations = service.Operations
	Handlers   = handlers.Handlers
)

func NewRepo(q *sqlc.Queries) *Repo { return repo.New(q) }

func New(d Deps) *Operations { return service.New(d) }

func NewHandlers(ops *Operations) *Handlers { return handlers.New(ops) }
