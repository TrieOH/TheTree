package repos

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"lib/telemetry"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.APIKey) (*models.APIKey, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).CreateApiKey(ctx, sqlc.CreateApiKeyParams{
		SubjectID:     toCreate.SubjectID,
		Name:          toCreate.Name,
		DisplayPrefix: toCreate.DisplayPrefix,
		KeyHash:       toCreate.KeyHash,
		CreatedBy:     toCreate.CreatedBy,
		ExpiresAt:     toCreate.ExpiresAt,
	})
	return new(mapAPIKey(row)), repo.dbe(err)
}
