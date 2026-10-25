package repo

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) SetFeeBPS(ctx context.Context, walletID uuid.UUID, feeBPS int) error {
	err := database.Queries(ctx, repo.q).SetWalletFeeBPS(ctx, sqlc.SetWalletFeeBPSParams{
		FeeBps: feeBPS,
		ID:     walletID,
	})
	return repo.dbe(err)
}
