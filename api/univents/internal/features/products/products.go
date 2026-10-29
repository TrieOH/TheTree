// Package products is the products feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package products

import (
	"univents/internal/features/products/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
