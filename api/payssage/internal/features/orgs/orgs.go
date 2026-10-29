// Package orgs is the orgs feature: its repo and handlers live here, its
// operations in service/, and Build wires the feature from here.
package orgs

import (
	"payssage/internal/features/orgs/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
