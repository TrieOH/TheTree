// Package webhook_events is the webhook events feature: the layers live in the
// subpackages, and Build wires the feature from here.
package webhook_events

import (
	"payssage/internal/features/webhook_events/handlers"
	"payssage/internal/features/webhook_events/repo"
	"payssage/internal/features/webhook_events/service"
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
