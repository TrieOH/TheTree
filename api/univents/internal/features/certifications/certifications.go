// Package certifications is the certifications feature: the layers live in the
// subpackages, and Build wires the feature from here.
package certifications

import (
	"univents/internal/features/certifications/handlers"
	"univents/internal/features/certifications/repo"
	"univents/internal/features/certifications/service"
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
