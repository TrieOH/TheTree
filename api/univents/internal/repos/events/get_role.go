package events

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetRole(ctx context.Context, actorID, eventID uuid.UUID) (models.EventMemberRole, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.GetRole")
	defer span.End()

	member, err := repo.GetMember(ctx, eventID, actorID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
