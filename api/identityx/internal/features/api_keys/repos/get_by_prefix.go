package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *Repo) GetByPrefix(ctx context.Context, prefix string) (*models.APIKey, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByPrefix")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetApiKeyByPrefix(ctx, prefix)
	return new(mapAPIKey(row)), repo.dbe(err)
}
