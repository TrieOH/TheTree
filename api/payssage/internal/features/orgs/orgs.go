// Package orgs is the orgs feature: the layers live in the
// subpackages, and Build wires the feature from here.
package orgs

import (
	"payssage/internal/features/orgs/handlers"
	"payssage/internal/features/orgs/repo"
	"payssage/internal/features/orgs/service"
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
