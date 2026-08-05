package wallets

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) SetFeeBPS(ctx context.Context, walletID uuid.UUID, feeBPS int) error {
	ctx, span := telemetry.StartSpan(ctx, "SetFeeBPS")
	defer span.End()
	err := database.Queries(ctx, repo.q).SetWalletFeeBPS(ctx, sqlc.SetWalletFeeBPSParams{
		FeeBps: feeBPS,
		ID:     walletID,
	})
	return repo.dbe(err)
}
