// Package webhook_events is the webhook events feature: its repo and
// handlers live here, its operations in service/, and Build wires the
// feature from here.
package webhook_events

import (
	"payssage/internal/features/webhook_events/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
