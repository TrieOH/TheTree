package repo

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *Repo) CreateState(ctx context.Context, state models.OAuthLoginState) (*models.OAuthLoginState, error) {
	row, err := database.Queries(ctx, repo.q).CreateOAuthLoginState(ctx, sqlc.CreateOAuthLoginStateParams{
		State:     state.State,
		Provider:  string(state.Provider),
		ProjectID: state.ProjectID,
		ExpiresAt: state.ExpiresAt,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOAuthLoginState(row)), nil
}
