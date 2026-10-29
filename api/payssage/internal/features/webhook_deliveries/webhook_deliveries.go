// Package webhook_deliveries is the webhook deliveries feature: its repo
// and handlers live here, its operations in service/, and Build wires the
// feature from here.
package webhook_deliveries

import (
	"payssage/internal/features/webhook_deliveries/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
