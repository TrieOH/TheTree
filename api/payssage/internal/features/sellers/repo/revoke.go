package repo

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) Revoke(ctx context.Context, id uuid.UUID) error {
	return repo.dbe(database.Queries(ctx, repo.q).RevokeSeller(ctx, id))
}
