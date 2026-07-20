package repos

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Collector, error) {
	ctx, span := repo.tracer.Start(ctx, "CollectorRepo.GetByID")
	defer span.End()

	sqlcCollector, err := database.Queries(ctx, repo.q).GetCollectorByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapCollector(sqlcCollector)), nil
}
