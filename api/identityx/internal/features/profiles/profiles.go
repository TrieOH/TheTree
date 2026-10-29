// Package profiles is the profiles feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package profiles

import (
	"IdentityX/internal/features/profiles/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
