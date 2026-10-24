// Package projects is the projects feature: the layers live in the
// subpackages, and Build wires the feature from here.
package projects

import (
	"IdentityX/internal/features/projects/handlers"
	"IdentityX/internal/features/projects/repo"
	"IdentityX/internal/features/projects/service"
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
