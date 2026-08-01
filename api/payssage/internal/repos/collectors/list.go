package collectors

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"payssage/models"
)

func (repo *Repo) List(ctx context.Context) ([]models.Collector, error) {
	ctx, span := telemetry.StartSpan(ctx, "CollectorRepo.List")
	defer span.End()

	sqlcCollectors, err := database.Queries(ctx, repo.q).ListCollectors(ctx)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcCollectors, mapCollector), nil
}
