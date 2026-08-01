package wallets

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) List(ctx context.Context, ownerID uuid.UUID) ([]models.Wallet, error) {
	ctx, span := telemetry.StartSpan(ctx, "List")
	defer span.End()
	sqlcWallets, err := database.Queries(ctx, repo.q).ListWallets(ctx, ownerID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcWallets, mapWallet), nil
}
