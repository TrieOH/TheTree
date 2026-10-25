package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Collector, error) {
	sqlcCollectors, err := database.Queries(ctx, repo.q).ListCollectorsByOrg(ctx, &orgID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcCollectors, mapCollector), nil
}
