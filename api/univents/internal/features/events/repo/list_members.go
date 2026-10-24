package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListEventMembers(ctx context.Context, eventID uuid.UUID) ([]models.EventMember, error) {
	members, err := database.Queries(ctx, repo.q).ListEventMembers(ctx, eventID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(members, mapEventMember), nil
}
