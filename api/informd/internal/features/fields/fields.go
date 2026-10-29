// Package fields is the fields feature: its repo and handlers live here,
// its operations in service/, and Build wires the feature from here.
package fields

import (
	"Informd/internal/features/fields/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
