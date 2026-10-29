// Package responses is the responses feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package responses

import (
	"Informd/internal/features/responses/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
