package repo

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (repo *Repo) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	sqlcWallet, err := database.Queries(ctx, repo.q).GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapWallet(sqlcWallet)), nil
}
