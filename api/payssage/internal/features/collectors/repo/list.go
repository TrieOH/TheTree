package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"
)

func (repo *Repo) List(ctx context.Context) ([]models.Collector, error) {
	sqlcCollectors, err := database.Queries(ctx, repo.q).ListCollectors(ctx)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcCollectors, mapCollector), nil
}
