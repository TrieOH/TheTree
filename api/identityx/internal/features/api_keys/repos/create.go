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
		SubjectID:     toCreate.SubjectID,
		Name:          toCreate.Name,
		DisplayPrefix: toCreate.DisplayPrefix,
		KeyHash:       toCreate.KeyHash,
		CreatedBy:     toCreate.CreatedBy,
		ExpiresAt:     toCreate.ExpiresAt,
	})
	return new(mapApiKey(row)), repo.dbe(err)
}
