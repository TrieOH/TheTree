package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"
)

func (repo *repo) CreateProduct(ctx context.Context, toCreate *models.Product) (*models.Product, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.CreateProduct")
	defer span.End()
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
