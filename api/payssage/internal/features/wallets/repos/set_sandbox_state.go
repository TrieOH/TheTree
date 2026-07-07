package repos

import (
	"context"
	"lib/database"
	"payssage/internal/database/sqlc"

	"github.com/google/uuid"
)

func (repo *repo) SetSandboxState(ctx context.Context, walletID uuid.UUID, state bool) error {
	ctx, span := repo.tracer.Start(ctx, "SetSandboxState")
	defer span.End()
	err := database.Queries(ctx, repo.q).SetWalletSandboxState(ctx, sqlc.SetWalletSandboxStateParams{
		Sandbox: state,
		ID:      walletID,
	})
	return repo.dbe(err)
}
