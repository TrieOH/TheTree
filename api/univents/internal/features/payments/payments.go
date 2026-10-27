// Package payments is the payments feature: the layers live in the
// subpackages, and Build wires the feature from here.
package payments

import (
	"univents/internal/features/payments/handlers"
	"univents/internal/features/payments/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
	Handlers   = handlers.Handlers
)

func New(d Deps) *Operations { return service.New(d) }

func NewHandlers(ops *Operations) *Handlers { return handlers.New(ops) }
