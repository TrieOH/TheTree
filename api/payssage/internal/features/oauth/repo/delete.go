package repo

import (
	"context"
	"lib/database"
)

func (repo *Repo) Delete(ctx context.Context, state string) error {
	err := database.Queries(ctx, repo.q).DeleteOAuthState(ctx, state)
	return repo.dbe(err)
}
