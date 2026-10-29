// Package webhooks is the webhooks feature: its handlers live here, its
// operations in service/, and Build wires the feature from here.
package webhooks

import (
	"univents/internal/features/webhooks/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
