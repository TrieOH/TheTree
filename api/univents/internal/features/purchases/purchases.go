// Package purchases is the purchases feature: the layers live in the
// subpackages, and Build wires the feature from here.
package purchases

import (
	"univents/internal/features/purchases/handlers"
	"univents/internal/features/purchases/repo"
	"univents/internal/features/purchases/service"
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
