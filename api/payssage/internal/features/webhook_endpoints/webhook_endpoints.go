// Package webhook_endpoints is the webhook endpoints feature: the layers live in the
// subpackages, and Build wires the feature from here.
package webhook_endpoints

import (
	"payssage/internal/features/webhook_endpoints/handlers"
	"payssage/internal/features/webhook_endpoints/repo"
	"payssage/internal/features/webhook_endpoints/service"
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
