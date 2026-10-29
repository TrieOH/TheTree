// Package programs is the programs feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package programs

import (
	"univents/internal/features/programs/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
