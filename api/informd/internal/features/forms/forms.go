// Package forms is the forms feature: the layers live in the
// subpackages, and Build wires the feature from here.
package forms

import (
	"Informd/internal/features/forms/handlers"
	"Informd/internal/features/forms/repo"
	"Informd/internal/features/forms/service"
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
