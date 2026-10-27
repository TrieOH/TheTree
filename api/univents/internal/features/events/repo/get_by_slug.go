package repo

import (
	"context"
	"lib/database"
	"univents/models"
)

func (repo *Repo) GetBySlug(ctx context.Context, slug string) (*models.Event, error) {
	event, err := database.Queries(ctx, repo.q).GetEventBySlug(ctx, slug)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapEvent(event)), nil
}
