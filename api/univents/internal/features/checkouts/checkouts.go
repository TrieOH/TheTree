// Package checkouts is the checkouts feature: the layers live in the
// subpackages, and Build wires the feature from here.
package checkouts

import (
	"univents/ports"

	idx "sdk/identityx"

	"univents/internal/features/checkouts/handlers"
	"univents/internal/features/checkouts/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
	Handlers   = handlers.Handlers
)

func New(d Deps) *Operations { return service.New(d) }

func NewHandlers(ops *Operations) *Handlers { return handlers.New(ops) }

// NewSDKActorResolver resolves attendees through the IdentityX actor service.
func NewSDKActorResolver(actors *idx.ActorService) ports.ActorResolver {
	return service.NewSDKActorResolver(actors)
}
