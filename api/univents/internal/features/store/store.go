// Package store is the store feature: the layers live in the
// subpackages, and Build wires the feature from here.
package store

import (
	"univents/internal/features/store/handlers"
	"univents/internal/features/store/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
	Handlers   = handlers.Handlers
)

func New(d Deps) *Operations { return service.New(d) }

func NewHandlers(ops *Operations) *Handlers { return handlers.New(ops) }
