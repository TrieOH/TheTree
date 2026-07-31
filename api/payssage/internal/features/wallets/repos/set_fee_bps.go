package repos

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) SetFeeBPS(ctx context.Context, walletID uuid.UUID, feeBPS int) error {
	ctx, span := repo.tracer.Start(ctx, "SetFeeBPS")
	defer span.End()
	err := database.Queries(ctx, repo.q).SetWalletFeeBPS(ctx, sqlc.SetWalletFeeBPSParams{
		FeeBps: feeBPS,
		ID:     walletID,
	})
	return repo.dbe(err)
}
