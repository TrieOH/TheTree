package repo

import (
	"context"
	"lib/database"
)

func (repo *Repo) ExpireStaleRequests(ctx context.Context) error {
	err := database.Queries(ctx, repo.q).ExpireStaleSignatureRequests(ctx)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
