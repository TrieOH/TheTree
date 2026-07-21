package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *repo) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Collector, error) {
	ctx, span := repo.tracer.Start(ctx, "CollectorRepo.ListByOwner")
	defer span.End()

	sqlcCollectors, err := database.Queries(ctx, repo.q).ListCollectorsByOwner(ctx, ownerID)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return xslices.MapSlice(sqlcCollectors, mapCollector), nil
}
