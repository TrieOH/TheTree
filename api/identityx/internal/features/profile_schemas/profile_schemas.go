// Package profile_schemas is the profile schemas feature: its repo and
// handlers live here, its operations in service/, and Build wires the
// feature from here.
package profile_schemas

import (
	"IdentityX/internal/features/profile_schemas/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
