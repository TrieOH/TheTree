package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) GetByPrefix(ctx context.Context, prefix string) (*models.ApiKey, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByPrefix")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetApiKeyByPrefix(ctx, prefix)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapApiKey(row)), nil
}
