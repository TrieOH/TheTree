package oauth_providers

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ProjectOAuthProvidersRepo = (*Repo)(nil)
var _ ports.OAuthLoginStatesRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("project oauth provider"),
	}
}

func mapProjectOAuthProvider(src sqlc.ProjectOauthProvider) models.ProjectOAuthProviders {
	return models.ProjectOAuthProviders{
		ID:                    src.ID,
		ProjectID:             src.ProjectID,
		Provider:              models.OAuthProvider(src.Provider),
		ClientID:              src.ClientID,
		EncryptedClientSecret: src.EncryptedClientSecret,
		CallbackURL:           src.CallbackUrl,
		Scopes:                src.Scopes,
		Enabled:               src.Enabled,
		CreatedAt:             src.CreatedAt,
		UpdatedAt:             src.UpdatedAt,
	}
}

func mapOAuthLoginState(src sqlc.OauthLoginState) models.OAuthLoginState {
	return models.OAuthLoginState{
		ID:        src.ID,
		State:     src.State,
		Provider:  models.OAuthProvider(src.Provider),
		ProjectID: src.ProjectID,
		CreatedAt: src.CreatedAt,
		ExpiresAt: src.ExpiresAt,
	}
}
