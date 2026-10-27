package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) RemoveEventMember(ctx context.Context, eventID, userID uuid.UUID) error {
	err := database.Queries(ctx, repo.q).RemoveEventMember(ctx, sqlc.RemoveEventMemberParams{
		EventID: eventID,
		UserID:  userID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
