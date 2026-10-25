package repo

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) SetSandboxState(ctx context.Context, walletID uuid.UUID, state bool) error {
	err := database.Queries(ctx, repo.q).SetWalletSandboxState(ctx, sqlc.SetWalletSandboxStateParams{
		Sandbox: state,
		ID:      walletID,
	})
	return repo.dbe(err)
}
