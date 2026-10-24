// Package badges is the badges feature: the layers live in the
// subpackages, and Build wires the feature from here.
package badges

import (
	"univents/internal/features/badges/handlers"
	"univents/internal/features/badges/repo"
	"univents/internal/features/badges/service"
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
