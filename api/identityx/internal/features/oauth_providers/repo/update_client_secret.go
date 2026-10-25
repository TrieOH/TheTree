package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) UpdateClientSecret(ctx context.Context, id uuid.UUID, encryptedSecret string) (*models.ProjectOAuthProviders, error) {
	row, err := database.Queries(ctx, repo.q).UpdateProjectOAuthProviderClientSecret(ctx, sqlc.UpdateProjectOAuthProviderClientSecretParams{
		ID:                    id,
		EncryptedClientSecret: encryptedSecret,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
