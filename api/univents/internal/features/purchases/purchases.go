// Package purchases is the purchases feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package purchases

import (
	"univents/internal/features/purchases/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
