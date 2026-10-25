// Package organizations is the organizations feature: the layers live in the
// subpackages, and Build wires the feature from here.
package organizations

import (
	"IdentityX/internal/features/organizations/handlers"
	"IdentityX/internal/features/organizations/repo"
	"IdentityX/internal/features/organizations/service"
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
