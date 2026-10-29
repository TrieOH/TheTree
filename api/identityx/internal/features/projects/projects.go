// Package projects is the projects feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package projects

import (
	"IdentityX/internal/features/projects/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
