package authz

import (
	"context"

	libauthz "lib/authz"
	"univents/ports"

	"github.com/google/uuid"
)

type Service struct {
	events ports.EventRepo
}

func New(events ports.EventRepo) *Service {
	return &Service{events: events}
}

func (s *Service) CheckEvent(ctx context.Context, actorID, eventID uuid.UUID, minRole libauthz.Role) error {
	role, err := s.events.GetRole(ctx, actorID, eventID)
	if err != nil {
		return libauthz.ForbiddenIfNotFound(err)
	}
	return libauthz.Min(role, minRole)
}
