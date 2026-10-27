package repo

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) DeleteState(ctx context.Context, id uuid.UUID) error {
	err := database.Queries(ctx, repo.q).DeleteOAuthLoginState(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
