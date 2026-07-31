package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()
	sqlcWallet, err := database.Queries(ctx, repo.q).GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapWallet(sqlcWallet)), nil
}
