// Package organizations is the organizations feature: its repo and
// handlers live here, its operations in service/, and Build wires the
// feature from here.
package organizations

import (
	"IdentityX/internal/features/organizations/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
