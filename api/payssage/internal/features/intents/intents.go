// Package intents is the intents feature: the layers live in the
// subpackages, and Build wires the feature from here.
package intents

import (
	"payssage/internal/features/intents/handlers"
	"payssage/internal/features/intents/repo"
	"payssage/internal/features/intents/service"
	"payssage/internal/sqlc"
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
