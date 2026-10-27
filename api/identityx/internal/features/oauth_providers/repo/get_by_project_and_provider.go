package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) GetByProjectAndProvider(ctx context.Context, projectID uuid.UUID, provider models.OAuthProvider) (*models.ProjectOAuthProviders, error) {
	row, err := database.Queries(ctx, repo.q).GetProjectOAuthProviderByProjectAndProvider(ctx, sqlc.GetProjectOAuthProviderByProjectAndProviderParams{
		ProjectID: projectID,
		Provider:  string(provider),
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
