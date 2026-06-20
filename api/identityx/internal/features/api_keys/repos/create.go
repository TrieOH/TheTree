package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) Create(ctx context.Context, toCreate models.ApiKey) (*models.ApiKey, error) {
	ctx, span := database.Span(ctx, repo.tracer, "Create")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).CreateApiKey(ctx, sqlc.CreateApiKeyParams{
		ActorID:   toCreate.ActorID,
		ProjectID: toCreate.ProjectID,
		Name:      toCreate.Name,
		KeyPrefix: toCreate.KeyPrefix,
		KeyHash:   toCreate.KeyHash,
		ExpiresAt: toCreate.ExpiresAt,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapApiKey(row)), nil
}
