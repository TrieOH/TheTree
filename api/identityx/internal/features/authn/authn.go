// Package authn is the authn feature: its repo and handlers live here, its
// operations in service/, and Build wires the feature from here.
package authn

import (
	"IdentityX/internal/features/authn/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }

// NewHandlers takes the setup state: until setup is done, only setup works.
