// Package authn is the authn feature: the layers live in the
// subpackages, and Build wires the feature from here.
package authn

import (
	"IdentityX/internal/features/authn/handlers"
	"IdentityX/internal/features/authn/repo"
	"IdentityX/internal/features/authn/service"
	"IdentityX/internal/setup"
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

// NewHandlers takes the setup state: until setup is done, only setup works.
func NewHandlers(ops *Operations, st *setup.State) *Handlers { return handlers.New(ops, st) }
