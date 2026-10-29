package oauth_providers

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

var _ ports.ProjectOAuthProvidersRepo = (*Repo)(nil)
var _ ports.OAuthLoginStatesRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrOAuthProviderNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

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

func (repo *Repo) Create(ctx context.Context, toCreate models.ProjectOAuthProviders) (*models.ProjectOAuthProviders, error) {
	row, err := repo.db(ctx).CreateProjectOAuthProvider(ctx, sqlc.CreateProjectOAuthProviderParams{
		ProjectID:             toCreate.ProjectID,
		Provider:              string(toCreate.Provider),
		ClientID:              toCreate.ClientID,
		EncryptedClientSecret: toCreate.EncryptedClientSecret,
		CallbackUrl:           toCreate.CallbackURL,
	})
	return database.One(row, err, repo.dbe, mapProjectOAuthProvider)
}

func (repo *Repo) CreateState(ctx context.Context, state models.OAuthLoginState) (*models.OAuthLoginState, error) {
	row, err := repo.db(ctx).CreateOAuthLoginState(ctx, sqlc.CreateOAuthLoginStateParams{
		State:     state.State,
		Provider:  string(state.Provider),
		ProjectID: state.ProjectID,
		ExpiresAt: state.ExpiresAt,
	})
	return database.One(row, err, repo.dbe, mapOAuthLoginState)
}

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).DeleteProjectOAuthProvider(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) DeleteState(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).DeleteOAuthLoginState(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.ProjectOAuthProviders, error) {
	row, err := repo.db(ctx).GetProjectOAuthProvider(ctx, id)
	return database.One(row, err, repo.dbe, mapProjectOAuthProvider)
}

func (repo *Repo) GetByProjectAndProvider(ctx context.Context, projectID uuid.UUID, provider models.OAuthProvider) (*models.ProjectOAuthProviders, error) {
	row, err := repo.db(ctx).GetProjectOAuthProviderByProjectAndProvider(ctx, sqlc.GetProjectOAuthProviderByProjectAndProviderParams{
		ProjectID: projectID,
		Provider:  string(provider),
	})
	return database.One(row, err, repo.dbe, mapProjectOAuthProvider)
}

func (repo *Repo) GetByState(ctx context.Context, state string) (*models.OAuthLoginState, error) {
	row, err := repo.db(ctx).GetOAuthLoginState(ctx, state)
	return database.One(row, err, repo.dbe, mapOAuthLoginState)
}

func (repo *Repo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.ProjectOAuthProviders, error) {
	rows, err := repo.db(ctx).ListProjectOAuthProviders(ctx, projectID)
	return database.Many(rows, err, repo.dbe, mapProjectOAuthProvider)
}

func (repo *Repo) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) (*models.ProjectOAuthProviders, error) {
	row, err := repo.db(ctx).SetProjectOAuthProviderEnabled(ctx, sqlc.SetProjectOAuthProviderEnabledParams{
		ID:      id,
		Enabled: enabled,
	})
	return database.One(row, err, repo.dbe, mapProjectOAuthProvider)
}

func (repo *Repo) UpdateCallbackURL(ctx context.Context, id uuid.UUID, callbackURL string) (*models.ProjectOAuthProviders, error) {
	row, err := repo.db(ctx).UpdateProjectOAuthProviderCallbackURL(ctx, sqlc.UpdateProjectOAuthProviderCallbackURLParams{
		ID:          id,
		CallbackUrl: callbackURL,
	})
	return database.One(row, err, repo.dbe, mapProjectOAuthProvider)
}

func (repo *Repo) UpdateClientID(ctx context.Context, id uuid.UUID, clientID string) (*models.ProjectOAuthProviders, error) {
	row, err := repo.db(ctx).UpdateProjectOAuthProviderClientID(ctx, sqlc.UpdateProjectOAuthProviderClientIDParams{
		ID:       id,
		ClientID: clientID,
	})
	return database.One(row, err, repo.dbe, mapProjectOAuthProvider)
}

func (repo *Repo) UpdateClientSecret(ctx context.Context, id uuid.UUID, encryptedSecret string) (*models.ProjectOAuthProviders, error) {
	row, err := repo.db(ctx).UpdateProjectOAuthProviderClientSecret(ctx, sqlc.UpdateProjectOAuthProviderClientSecretParams{
		ID:                    id,
		EncryptedClientSecret: encryptedSecret,
	})
	return database.One(row, err, repo.dbe, mapProjectOAuthProvider)
}
