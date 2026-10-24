package repo

import (
	"context"

	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := database.Queries(ctx, repo.q).DeleteField(ctx, id)
	if err != nil {
		return repo.dbe(err)
	}
	return nil
}
