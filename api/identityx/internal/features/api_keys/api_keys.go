// Package api_keys is the api keys feature: the layers live in the
// subpackages, and Build wires the feature from here.
package api_keys

import (
	"IdentityX/internal/features/api_keys/handlers"
	"IdentityX/internal/features/api_keys/repo"
	"IdentityX/internal/features/api_keys/service"
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
