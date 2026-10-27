package repo

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ExternalIdentitiesRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrExternalIdentityNotFound),
	}
}

func mapExternalIdentity(src sqlc.ActorExternalIdentity) models.ActorExternalIdentities {
	return models.ActorExternalIdentities{
		ID:                    src.ID,
		ActorID:               src.ActorID,
		Provider:              models.OAuthProvider(src.Provider),
		Subject:               src.Subject,
		Email:                 src.Email,
		EncryptedAccessToken:  src.EncryptedAccessToken,
		EncryptedRefreshToken: src.EncryptedRefreshToken,
		TokenExpiresAt:        src.TokenExpiresAt,
		CreatedAt:             src.CreatedAt,
		UpdatedAt:             src.UpdatedAt,
	}
}
