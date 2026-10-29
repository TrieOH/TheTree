// Package api_keys is the api keys feature: its repo and handlers live
// here, its operations in service/, and Build wires the feature from here.
package api_keys

import (
	"IdentityX/internal/features/api_keys/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
