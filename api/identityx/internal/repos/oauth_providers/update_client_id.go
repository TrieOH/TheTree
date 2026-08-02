package oauth_providers

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) UpdateClientID(ctx context.Context, id uuid.UUID, clientID string) (*models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpdateClientID")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).UpdateProjectOAuthProviderClientID(ctx, sqlc.UpdateProjectOAuthProviderClientIDParams{
		ID:       id,
		ClientID: clientID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
