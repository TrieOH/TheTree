// Package steps is the steps feature: the layers live in the
// subpackages, and Build wires the feature from here.
package steps

import (
	"Informd/internal/features/steps/handlers"
	"Informd/internal/features/steps/repo"
	"Informd/internal/features/steps/service"
	"Informd/internal/sqlc"
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
