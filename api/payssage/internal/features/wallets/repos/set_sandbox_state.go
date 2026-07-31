package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/internal/sqlc"

	"github.com/google/uuid"
)

func (repo *Repo) SetSandboxState(ctx context.Context, walletID uuid.UUID, state bool) error {
	ctx, span := telemetry.StartSpan(ctx, "SetSandboxState")
	defer span.End()
	err := database.Queries(ctx, repo.q).SetWalletSandboxState(ctx, sqlc.SetWalletSandboxStateParams{
		Sandbox: state,
		ID:      walletID,
	})
	return repo.dbe(err)
}
