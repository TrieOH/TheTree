package oauth_providers

import (
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) GetByState(ctx context.Context, state string) (*models.OAuthLoginState, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByState")
	defer span.End()

	row, err := database.Queries(ctx, repo.q).GetOAuthLoginState(ctx, state)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOAuthLoginState(row)), nil
}
