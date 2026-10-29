// Package actors is the actors feature: its repo and handlers live here,
// its operations in service/, and Build wires the feature from here.
package actors

import (
	"IdentityX/internal/features/actors/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
