// Package programs is the programs feature: the layers live in the
// subpackages, and Build wires the feature from here.
package programs

import (
	"univents/internal/features/programs/handlers"
	"univents/internal/features/programs/repo"
	"univents/internal/features/programs/service"
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
