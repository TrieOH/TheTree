package repos

import (
	"context"
	"lib/database"
	"payssage/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) BindCollector(ctx context.Context, walletID, collectorID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "WalletRepo.BindCollector")
	defer span.End()

	return repo.dbe(database.Queries(ctx, repo.q).BindCollector(ctx, sqlc.BindCollectorParams{
		ID:          walletID,
		CollectorID: &collectorID,
	}))
}
