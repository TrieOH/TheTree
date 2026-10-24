// Package capabilities is the capabilities feature: the layers live in the
// subpackages, and Build wires the feature from here.
package capabilities

import (
	"IdentityX/internal/features/capabilities/handlers"
	"IdentityX/internal/features/capabilities/repo"
	"IdentityX/internal/features/capabilities/service"
	"IdentityX/internal/sqlc"
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
