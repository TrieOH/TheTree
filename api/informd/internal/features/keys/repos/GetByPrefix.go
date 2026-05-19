package repos

import (
	"Informd/models"
	"context"
	"lib/database"
	"lib/xslices"
)

func (repo *repo) GetByPrefix(ctx context.Context, prefix string) ([]models.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByPrefix")
	defer span.End()
	sqlcApiKeys, err := database.Queries(ctx, repo.q).GetAPIKeyByPrefix(ctx, prefix)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcApiKeys, mapApiKey), nil
}
