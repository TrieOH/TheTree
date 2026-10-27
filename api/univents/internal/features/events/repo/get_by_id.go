package repo

import (
	"context"
	"lib/database"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	event, err := database.Queries(ctx, repo.q).GetEventByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(event)), nil
}
