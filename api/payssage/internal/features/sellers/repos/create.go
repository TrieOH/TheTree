package repos

import (
	"context"
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
)

func (repo *Repo) Create(ctx context.Context, toCreate models.Seller) (*models.Seller, error) {
	ctx, span := repo.tracer.Start(ctx, "SellerRepo.Create")
	defer span.End()

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
