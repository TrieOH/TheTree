package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Collector, error) {
	sqlcCollectors, err := database.Queries(ctx, repo.q).ListCollectorsByOwner(ctx, ownerID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcCollectors, mapCollector), nil
}
