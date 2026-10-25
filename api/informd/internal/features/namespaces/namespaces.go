// Package namespaces is the namespaces feature: the layers live in the
// subpackages, and Build wires the feature from here.
package namespaces

import (
	"Informd/internal/features/namespaces/handlers"
	"Informd/internal/features/namespaces/repo"
	"Informd/internal/features/namespaces/service"
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
