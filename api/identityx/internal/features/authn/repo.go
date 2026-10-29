package authn

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ExternalIdentitiesRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrExternalIdentityNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

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

func (repo *Repo) Create(ctx context.Context, identity models.ActorExternalIdentities) (*models.ActorExternalIdentities, error) {
	row, err := repo.db(ctx).CreateExternalIdentity(ctx, sqlc.CreateExternalIdentityParams{
		ActorID:               identity.ActorID,
		Provider:              string(identity.Provider),
		Subject:               identity.Subject,
		Email:                 identity.Email,
		EncryptedAccessToken:  identity.EncryptedAccessToken,
		EncryptedRefreshToken: identity.EncryptedRefreshToken,
		TokenExpiresAt:        identity.TokenExpiresAt,
	})
	return database.One(row, err, repo.dbe, mapExternalIdentity)
}

func (repo *Repo) GetByProviderAndSubject(ctx context.Context, provider, subject string, projectID *uuid.UUID) (*models.ActorExternalIdentities, error) {
	row, err := repo.db(ctx).GetExternalIdentityByProviderAndSubject(ctx, sqlc.GetExternalIdentityByProviderAndSubjectParams{
		Provider:  provider,
		Subject:   subject,
		ProjectID: projectID,
	})
	return database.One(row, err, repo.dbe, mapExternalIdentity)
}

func (repo *Repo) UpdateTokens(ctx context.Context, identity models.ActorExternalIdentities, projectID *uuid.UUID) (*models.ActorExternalIdentities, error) {
	row, err := repo.db(ctx).UpdateExternalIdentityTokens(ctx, sqlc.UpdateExternalIdentityTokensParams{
		Provider:              string(identity.Provider),
		Subject:               identity.Subject,
		EncryptedAccessToken:  identity.EncryptedAccessToken,
		EncryptedRefreshToken: identity.EncryptedRefreshToken,
		TokenExpiresAt:        identity.TokenExpiresAt,
		ProjectID:             projectID,
	})
	return database.One(row, err, repo.dbe, mapExternalIdentity)
}
