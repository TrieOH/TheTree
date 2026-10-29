// Package steps is the steps feature: its repo and handlers live here, its
// operations in service/, and Build wires the feature from here.
package steps

import (
	"Informd/internal/features/steps/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
