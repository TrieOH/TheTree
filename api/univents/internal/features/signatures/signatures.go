// Package signatures is the signatures feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package signatures

import (
	"univents/internal/features/signatures/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
