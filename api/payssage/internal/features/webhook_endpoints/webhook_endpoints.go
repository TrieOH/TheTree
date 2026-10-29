// Package webhook_endpoints is the webhook endpoints feature: its repo and
// handlers live here, its operations in service/, and Build wires the
// feature from here.
package webhook_endpoints

import (
	"payssage/internal/features/webhook_endpoints/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
