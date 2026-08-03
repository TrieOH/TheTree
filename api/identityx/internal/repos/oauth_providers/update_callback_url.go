package oauth_providers

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) UpdateCallbackURL(ctx context.Context, id uuid.UUID, callbackURL string) (*models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpdateCallbackURL")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).UpdateProjectOAuthProviderCallbackURL(ctx, sqlc.UpdateProjectOAuthProviderCallbackURLParams{
		ID:          id,
		CallbackUrl: callbackURL,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
