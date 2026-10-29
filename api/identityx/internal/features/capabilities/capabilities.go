// Package capabilities is the capabilities feature: its repo and handlers
// live here, its operations in service/, and Build wires the feature from
// here.
package capabilities

import (
	"IdentityX/internal/features/capabilities/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
