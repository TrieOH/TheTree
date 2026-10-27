// Package sellers is the sellers feature: the layers live in the
// subpackages, and Build wires the feature from here.
package sellers

import (
	"payssage/internal/features/sellers/handlers"
	"payssage/internal/features/sellers/repo"
	"payssage/internal/features/sellers/service"
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
