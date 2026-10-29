// Package oauth is the oauth feature: its repo and handlers live here, its
// operations in service/, and Build wires the feature from here.
package oauth

import (
	"payssage/internal/features/oauth/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
