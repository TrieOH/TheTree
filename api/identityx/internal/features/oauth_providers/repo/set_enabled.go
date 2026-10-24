package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) (*models.ProjectOAuthProviders, error) {
	row, err := database.Queries(ctx, repo.q).SetProjectOAuthProviderEnabled(ctx, sqlc.SetProjectOAuthProviderEnabledParams{
		ID:      id,
		Enabled: enabled,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
