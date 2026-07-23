package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) RemoveEventMember(ctx context.Context, eventID, userID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.RemoveEventMember")
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
