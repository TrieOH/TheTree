package oauth_providers

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) GetByProjectAndProvider(ctx context.Context, projectID uuid.UUID, provider models.OAuthProvider) (*models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByProjectAndProvider")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).GetProjectOAuthProviderByProjectAndProvider(ctx, sqlc.GetProjectOAuthProviderByProjectAndProviderParams{
		ProjectID: projectID,
		Provider:  string(provider),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
