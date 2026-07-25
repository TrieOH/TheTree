package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) CreateProduct(ctx context.Context, toCreate *models.Product) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsRepo.CreateProduct")
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
