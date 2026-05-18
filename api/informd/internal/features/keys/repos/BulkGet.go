package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (repo *repo) BulkGet(ctx context.Context, ids []uuid.UUID) ([]models.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "BulkGet")
	defer span.End()
	sqlcKeys, err := database.Queries(ctx, repo.q).BulkGetAPIKeys(ctx, ids)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcKeys, mapApiKey), nil
}
