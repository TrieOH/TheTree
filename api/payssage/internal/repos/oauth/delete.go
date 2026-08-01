package oauth

import (
	"context"
	"lib/database"
	"lib/telemetry"
)

func (repo *Repo) Delete(ctx context.Context, state string) error {
	ctx, span := telemetry.StartSpan(ctx, "Delete")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteOAuthState(ctx, state)
	return repo.dbe(err)
}
