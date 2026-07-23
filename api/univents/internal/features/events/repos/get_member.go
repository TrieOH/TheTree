package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) GetMember(ctx context.Context, eventID, userID uuid.UUID) (*models.EventMember, error) {
	ctx, span := repo.tracer.Start(ctx, "EventsRepo.GetEventMember")
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
