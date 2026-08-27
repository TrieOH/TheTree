package checkouts

import (
	"context"

	idx "sdk/identityx"
	"univents/ports"

	"github.com/MintzyG/sdkkit"
	"github.com/google/uuid"
)

// sdkActors adapts the IdentityX SDK ActorService to the ActorResolver
// seam (ports.ActorResolver): the SDK surfaces "no account for this email"
// / "unknown actor" as a 404 APIError, which is mapped to
// ErrActorNotFound so the checkout logic can branch on "no account yet"
// (email-only gift) vs "verification failed" (error).
type sdkActors struct {
	actors *idx.ActorService
}

// NewSDKActorResolver wraps the IdentityX actor service for checkout's
// attendee resolution, gift claim, and the profile-badges gift claim.
// Lookups are scoped to the univents project (the client's project id).
func NewSDKActorResolver(actors *idx.ActorService) ports.ActorResolver {
	return &sdkActors{actors: actors}
}

func (s *sdkActors) GetByEmail(ctx context.Context, email string) (*idx.Actor, error) {
	actor, err := s.actors.GetByEmail(ctx, email)
	if err != nil {
		if sdkkit.IsNotFound(err) {
			return nil, ports.ErrActorNotFound
		}
		return nil, err
	}
	return actor, nil
}

func (s *sdkActors) GetByID(ctx context.Context, id uuid.UUID) (*idx.Actor, error) {
	actor, err := s.actors.GetByID(ctx, id)
	if err != nil {
		if sdkkit.IsNotFound(err) {
			return nil, ports.ErrActorNotFound
		}
		return nil, err
	}
	return actor, nil
}
