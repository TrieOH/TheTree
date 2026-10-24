package authz

import (
	"context"

	libauthz "lib/authz"
	"univents/ports"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

type Service struct {
	events ports.EventRepo
}

func New(events ports.EventRepo) *Service {
	return &Service{events: events}
}

// Event resolves the caller from the identity the auth middleware wrote and
// checks they hold at least minRole on the event. It returns the caller, so
// operations don't look it up again.
func (s *Service) Event(ctx context.Context, eventID uuid.UUID, minRole libauthz.Role) (*idx.Identity, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	return ident, s.EventFor(ctx, ident.Sub.ID, eventID, minRole)
}

// EventFor checks a given actor instead of the caller: an attendee being
// admitted to a staff-only program, or an operation handed its actor.
func (s *Service) EventFor(ctx context.Context, actorID, eventID uuid.UUID, minRole libauthz.Role) error {
	role, err := s.events.GetRole(ctx, actorID, eventID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}
