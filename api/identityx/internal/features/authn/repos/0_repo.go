package repos

import (
	sqlc2 "IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc2.Queries
	dbe database.ErrorHandler
}

var _ ports.ExternalIdentitiesRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("external actor identity"),
	}
}

func mapExternalIdentity(src sqlc2.ActorExternalIdentity) models.ActorExternalIdentities {
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
