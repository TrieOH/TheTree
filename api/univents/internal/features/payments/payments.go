// Package payments is the payments feature: its handlers live here, its
// operations in service/, and Build wires the feature from here.
package payments

import (
	"univents/internal/features/payments/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
