// Package ticket_types is the ticket types feature: the layers live in the
// subpackages, and Build wires the feature from here.
package ticket_types

import (
	"univents/internal/features/ticket_types/handlers"
	"univents/internal/features/ticket_types/repo"
	"univents/internal/features/ticket_types/service"
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
