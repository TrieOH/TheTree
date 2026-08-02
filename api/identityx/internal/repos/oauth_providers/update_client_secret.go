package oauth_providers

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) UpdateClientSecret(ctx context.Context, id uuid.UUID, encryptedSecret string) (*models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpdateClientSecret")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).UpdateProjectOAuthProviderClientSecret(ctx, sqlc.UpdateProjectOAuthProviderClientSecretParams{
		ID:                    id,
		EncryptedClientSecret: encryptedSecret,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectOAuthProvider(row)), nil
}
