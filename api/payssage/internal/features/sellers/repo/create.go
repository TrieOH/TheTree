package repo

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Seller) (*models.Seller, error) {
	sqlcSeller, err := database.Queries(ctx, repo.q).CreateSeller(ctx, sqlc.CreateSellerParams{
		WalletID:       toCreate.WalletID,
		Provider:       toCreate.Provider,
		ProviderUserID: toCreate.ProviderUserID,
		Credentials:    toCreate.Credentials,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}

	return new(mapSeller(sqlcSeller)), nil
}
