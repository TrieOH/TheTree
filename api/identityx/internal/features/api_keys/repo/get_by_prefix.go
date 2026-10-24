package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *Repo) GetByPrefix(ctx context.Context, prefix string) (*models.APIKey, error) {
	row, err := database.Queries(ctx, repo.q).GetApiKeyByPrefix(ctx, prefix)
	return new(mapAPIKey(row)), repo.dbe(err)
}
