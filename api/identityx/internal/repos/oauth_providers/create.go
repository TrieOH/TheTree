package oauth_providers

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.ProjectOAuthProviders) (*models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "Create")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).CreateProjectOAuthProvider(ctx, sqlc.CreateProjectOAuthProviderParams{
		ProjectID:             toCreate.ProjectID,
		Provider:              string(toCreate.Provider),
		ClientID:              toCreate.ClientID,
		EncryptedClientSecret: toCreate.EncryptedClientSecret,
		CallbackUrl:           toCreate.CallbackURL,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
