package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"
)

func (repo *Repo) ListPublic(ctx context.Context) ([]models.Event, error) {
	sqlcEvents, err := database.Queries(ctx, repo.q).ListPublicEvents(ctx)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcEvents, mapEvent), nil
}
