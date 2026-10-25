// Package ws is the ws feature: the layers live in the
// subpackages, and Build wires the feature from here.
package ws

import (
	"univents/internal/features/ws/handlers"
	"univents/internal/features/ws/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
	Handlers   = handlers.Handlers
)

func New(d Deps) *Operations { return service.New(d) }

func NewHandlers(ops *Operations) *Handlers { return handlers.New(ops) }
