package oauth

import (
	"context"
	"lib/database"
	"payssage/errs"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.OAuthStateRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrOAuthStateNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapState(src sqlc.OauthState) models.OAuthState {
	return models.OAuthState{
		State:            src.State,
		WalletID:         src.WalletID,
		OrganizationID:   src.OrganizationID,
		OwnerID:          src.OwnerID,
		Provider:         src.Provider,
		Flow:             models.OAuthFlow(src.Flow),
		FinalRedirectURL: src.FinalRedirectUrl,
		CreatedAt:        src.CreatedAt,
		ExpiresAt:        src.ExpiresAt,
	}
}

func (repo *Repo) Create(ctx context.Context, state models.OAuthState) (*models.OAuthState, error) {
	sqlcState, err := repo.db(ctx).CreateOAuthState(ctx, sqlc.CreateOAuthStateParams{
		State:            state.State,
		WalletID:         state.WalletID,
		OrganizationID:   state.OrganizationID,
		OwnerID:          state.OwnerID,
		Provider:         state.Provider,
		Flow:             state.Flow.String(),
		FinalRedirectUrl: state.FinalRedirectURL,
		ExpiresAt:        state.ExpiresAt,
	})
	return database.One(sqlcState, err, repo.dbe, mapState)
}

func (repo *Repo) Delete(ctx context.Context, state string) error {
	err := repo.db(ctx).DeleteOAuthState(ctx, state)
	return repo.dbe(err)
}

func (repo *Repo) Get(ctx context.Context, state string) (*models.OAuthState, error) {
	sqlcState, err := repo.db(ctx).GetOAuthState(ctx, state)
	return database.One(sqlcState, err, repo.dbe, mapState)
}
