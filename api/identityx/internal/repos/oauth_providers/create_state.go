package oauth_providers

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) CreateState(ctx context.Context, state models.OAuthLoginState) (*models.OAuthLoginState, error) {
	ctx, span := telemetry.StartSpan(ctx, "CreateState")
	defer span.End()

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
