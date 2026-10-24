package repo

import (
	"IdentityX/models"
	"context"
	"lib/database"
)

func (repo *Repo) GetByState(ctx context.Context, state string) (*models.OAuthLoginState, error) {
	row, err := database.Queries(ctx, repo.q).GetOAuthLoginState(ctx, state)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapOAuthLoginState(row)), nil
}
