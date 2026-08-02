package oauth_providers

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) (*models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "SetEnabled")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).SetProjectOAuthProviderEnabled(ctx, sqlc.SetProjectOAuthProviderEnabledParams{
		ID:      id,
		Enabled: enabled,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
