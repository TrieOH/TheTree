package events

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetMember(ctx context.Context, eventID, userID uuid.UUID) (*models.EventMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.GetEventMember")
	defer span.End()

	member, err := database.Queries(ctx, repo.q).GetEventMember(ctx, sqlc.GetEventMemberParams{
		EventID: eventID,
		UserID:  userID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapEventMember(member)), nil
}
