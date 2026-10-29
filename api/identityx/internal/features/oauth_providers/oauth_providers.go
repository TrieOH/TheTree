// Package oauth_providers is the oauth providers feature: its repo and
// handlers live here, its operations in service/, and Build wires the
// feature from here.
package oauth_providers

import (
	"IdentityX/internal/features/oauth_providers/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
