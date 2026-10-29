// Package checkouts is the checkouts feature: its handlers live here, its
// operations in service/, and Build wires the feature from here.
package checkouts

import (
	"univents/ports"

	idx "sdk/identityx"

	"univents/internal/features/checkouts/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }

// NewSDKActorResolver resolves attendees through the IdentityX actor service.
func NewSDKActorResolver(actors *idx.ActorService) ports.ActorResolver {
	return service.NewSDKActorResolver(actors)
}
