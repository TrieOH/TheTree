// Package sellers is the sellers feature: its repo and handlers live here,
// its operations in service/, and Build wires the feature from here.
package sellers

import (
	"payssage/internal/features/sellers/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
