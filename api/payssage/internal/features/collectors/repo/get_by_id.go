package repo

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Collector, error) {
	sqlcCollector, err := database.Queries(ctx, repo.q).GetCollectorByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapCollector(sqlcCollector)), nil
}
