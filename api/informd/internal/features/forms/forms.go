// Package forms is the forms feature: its repo and handlers live here, its
// operations in service/, and Build wires the feature from here.
package forms

import (
	"Informd/internal/features/forms/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
