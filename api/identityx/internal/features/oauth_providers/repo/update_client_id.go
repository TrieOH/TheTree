package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) UpdateClientID(ctx context.Context, id uuid.UUID, clientID string) (*models.ProjectOAuthProviders, error) {
	row, err := database.Queries(ctx, repo.q).UpdateProjectOAuthProviderClientID(ctx, sqlc.UpdateProjectOAuthProviderClientIDParams{
		ID:       id,
		ClientID: clientID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
