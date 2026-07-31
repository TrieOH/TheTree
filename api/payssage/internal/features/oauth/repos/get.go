package repos

import (
	"context"
	"lib/database"
	"payssage/models"
)

func (repo *Repo) Get(ctx context.Context, state string) (*models.OAuthState, error) {
	ctx, span := repo.tracer.Start(ctx, "Get")
	defer span.End()
	sqlcState, err := database.Queries(ctx, repo.q).GetOAuthState(ctx, state)
	return new(mapState(sqlcState)), repo.dbe(err)
}
