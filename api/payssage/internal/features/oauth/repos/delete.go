package repos

import (
	"context"
	"lib/database"
)

func (repo *Repo) Delete(ctx context.Context, state string) error {
	ctx, span := repo.tracer.Start(ctx, "Delete")
	defer span.End()
	err := database.Queries(ctx, repo.q).DeleteOAuthState(ctx, state)
	return repo.dbe(err)
}
