package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) CreateProduct(ctx context.Context, toCreate *models.Product) (*models.Product, error) {
	result, err := database.Queries(ctx, repo.q).CreateProduct(ctx, sqlc.CreateProductParams{
		EditionID:            toCreate.EditionID,
		VendorCode:           toCreate.VendorCode,
		RequiresRegistration: toCreate.RequiresRegistration,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProduct(result)), nil
}
