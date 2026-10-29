// Package namespaces is the namespaces feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package namespaces

import (
	"Informd/internal/features/namespaces/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
