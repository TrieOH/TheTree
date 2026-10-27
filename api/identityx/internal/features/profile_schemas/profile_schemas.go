// Package profile_schemas is the profile schemas feature: the layers live in the
// subpackages, and Build wires the feature from here.
package profile_schemas

import (
	"IdentityX/internal/features/profile_schemas/handlers"
	"IdentityX/internal/features/profile_schemas/repo"
	"IdentityX/internal/features/profile_schemas/service"
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
