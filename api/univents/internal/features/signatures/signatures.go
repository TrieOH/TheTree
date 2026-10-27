// Package signatures is the signatures feature: the layers live in the
// subpackages, and Build wires the feature from here.
package signatures

import (
	"univents/internal/features/signatures/handlers"
	"univents/internal/features/signatures/repo"
	"univents/internal/features/signatures/service"
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
