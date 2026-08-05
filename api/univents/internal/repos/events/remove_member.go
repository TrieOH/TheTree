package events

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveEventMember(ctx context.Context, eventID, userID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "EventsRepo.RemoveEventMember")
	defer span.End()

	err := database.Queries(ctx, repo.q).RemoveEventMember(ctx, sqlc.RemoveEventMemberParams{
		EventID: eventID,
		UserID:  userID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
