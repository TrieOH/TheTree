package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) BindCollector(ctx context.Context, walletID, collectorID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "WalletRepo.BindCollector")
	defer span.End()

	return repo.dbe(database.Queries(ctx, repo.q).BindCollector(ctx, sqlc.BindCollectorParams{
		ID:          walletID,
		CollectorID: &collectorID,
	}))
}
