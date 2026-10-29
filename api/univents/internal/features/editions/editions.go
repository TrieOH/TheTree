// Package editions is the editions feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package editions

import (
	"univents/internal/features/editions/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
