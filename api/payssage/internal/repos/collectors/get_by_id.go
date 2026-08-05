package collectors

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Collector, error) {
	ctx, span := telemetry.StartSpan(ctx, "CollectorRepo.GetByID")
	defer span.End()

	sqlcCollector, err := database.Queries(ctx, repo.q).GetCollectorByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapCollector(sqlcCollector)), nil
}
