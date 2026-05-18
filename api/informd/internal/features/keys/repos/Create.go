package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"context"
	"lib/database"
)

func (repo *repo) Create(ctx context.Context, toCreate models.APIKey) (*models.APIKey, error) {
	ctx, span := repo.tracer.Start(ctx, "Create")
	defer span.End()
	sqlcApiKey, err := database.Queries(ctx, repo.q).CreateAPIKey(ctx, sqlc.CreateAPIKeyParams{
		OwnerID:   toCreate.OwnerID,
		Name:      toCreate.Name,
		KeyHash:   toCreate.KeyHash,
		KeyPrefix: toCreate.KeyPrefix,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapApiKey(sqlcApiKey)), nil
}
