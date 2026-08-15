package authn

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) UpdateTokens(ctx context.Context, identity models.ActorExternalIdentities, projectID *uuid.UUID) (*models.ActorExternalIdentities, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpdateTokens")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).UpdateExternalIdentityTokens(ctx, sqlc.UpdateExternalIdentityTokensParams{
		Provider:              string(identity.Provider),
		Subject:               identity.Subject,
		EncryptedAccessToken:  identity.EncryptedAccessToken,
		EncryptedRefreshToken: identity.EncryptedRefreshToken,
		TokenExpiresAt:        identity.TokenExpiresAt,
		ProjectID:             projectID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapExternalIdentity(row)), nil
}
