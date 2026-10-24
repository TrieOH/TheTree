package repo

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) UnbindCollector(ctx context.Context, walletID uuid.UUID) error {
	return repo.dbe(database.Queries(ctx, repo.q).UnbindCollector(ctx, walletID))
}
