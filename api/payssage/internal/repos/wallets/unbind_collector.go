package wallets

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (repo *Repo) UnbindCollector(ctx context.Context, walletID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "WalletRepo.UnbindCollector")
	defer span.End()

	return repo.dbe(database.Queries(ctx, repo.q).UnbindCollector(ctx, walletID))
}
