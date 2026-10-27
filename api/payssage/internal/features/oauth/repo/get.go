package repo

import (
	"context"
	"lib/database"
	"payssage/models"
)

func (repo *Repo) Get(ctx context.Context, state string) (*models.OAuthState, error) {
	sqlcState, err := database.Queries(ctx, repo.q).GetOAuthState(ctx, state)
	return new(mapState(sqlcState)), repo.dbe(err)
}
