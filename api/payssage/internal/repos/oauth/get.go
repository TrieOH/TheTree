package oauth

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/models"
)

func (repo *Repo) Get(ctx context.Context, state string) (*models.OAuthState, error) {
	ctx, span := telemetry.StartSpan(ctx, "Get")
	defer span.End()
	sqlcState, err := database.Queries(ctx, repo.q).GetOAuthState(ctx, state)
	return new(mapState(sqlcState)), repo.dbe(err)
}
