package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *repo) UpdateTokens(ctx context.Context, identity models.ActorExternalIdentities) (*models.ActorExternalIdentities, error) {
	ctx, span := database.Span(ctx, repo.tracer, "UpdateTokens")
	defer span.End()
	row, err := database.Queries(ctx, repo.q).UpdateExternalIdentityTokens(ctx, sqlc.UpdateExternalIdentityTokensParams{
		Provider:              string(identity.Provider),
		Subject:               identity.Subject,
		EncryptedAccessToken:  identity.EncryptedAccessToken,
		EncryptedRefreshToken: identity.EncryptedRefreshToken,
		TokenExpiresAt:        identity.TokenExpiresAt,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapExternalIdentity(row)), nil
}
