package ports

import (
	"context"
	"errors"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

// ActorResolver is the IdentityX actor lookup seam: resolve an attendee
// email to its account in the univents identityx project (GetByEmail — gift
// buying), or a user's account email by id (GetByID — the gift claim, which
// matches gifted registrations by the account's own email). No account /
// unknown actor is reported as ErrActorNotFound. Satisfied by an adapter
// over the SDK's *idx.ActorService (checkouts.NewSDKActorResolver); faked
// in tests.
type ActorResolver interface {
	GetByEmail(ctx context.Context, email string) (*idx.Actor, error)
	GetByID(ctx context.Context, id uuid.UUID) (*idx.Actor, error)
}

// ErrActorNotFound is returned by ActorResolver when the email or id has no
// account in the univents IdentityX project.
var ErrActorNotFound = errors.New("identityx: no account for email")
