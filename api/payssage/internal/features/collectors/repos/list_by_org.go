package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *repo) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Collector, error) {
	ctx, span := repo.tracer.Start(ctx, "CollectorRepo.ListByOrg")
	defer span.End()

	sqlcCollectors, err := database.Queries(ctx, repo.q).ListCollectorsByOrg(ctx, &orgID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcCollectors, mapCollector), nil
}
