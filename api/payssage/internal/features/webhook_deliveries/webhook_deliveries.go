// Package webhook_deliveries is the webhook deliveries feature: the layers live in the
// subpackages, and Build wires the feature from here.
package webhook_deliveries

import (
	"payssage/internal/features/webhook_deliveries/handlers"
	"payssage/internal/features/webhook_deliveries/repo"
	"payssage/internal/features/webhook_deliveries/service"
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
