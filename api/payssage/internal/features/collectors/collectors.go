// Package collectors is the collectors feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package collectors

import (
	"payssage/internal/features/collectors/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
