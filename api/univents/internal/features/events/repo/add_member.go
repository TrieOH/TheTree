package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) AddEventMember(ctx context.Context, eventID, userID uuid.UUID, role models.EventMemberRole) (*models.EventMember, error) {
	member, err := database.Queries(ctx, repo.q).AddEventMember(ctx, sqlc.AddEventMemberParams{
		EventID: eventID,
		UserID:  userID,
		Role:    string(role),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapEventMember(member)), nil
}
