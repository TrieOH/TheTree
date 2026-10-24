// Package events is the events feature: the layers live in the
// subpackages, and Build wires the feature from here.
package events

import (
	"univents/internal/features/events/handlers"
	"univents/internal/features/events/repo"
	"univents/internal/features/events/service"
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
