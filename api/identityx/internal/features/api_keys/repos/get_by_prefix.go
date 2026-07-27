package repos

import (
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"
)

func (repo *Repo) GetByPrefix(ctx context.Context, prefix string) (*models.APIKey, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByPrefix")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).GetApiKeyByPrefix(ctx, prefix)
	return new(mapAPIKey(row)), repo.dbe(err)
}
