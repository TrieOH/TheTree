package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (repo *Repo) UnbindCollector(ctx context.Context, walletID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "WalletRepo.UnbindCollector")
	defer span.End()

	return repo.dbe(database.Queries(ctx, repo.q).UnbindCollector(ctx, walletID))
}
