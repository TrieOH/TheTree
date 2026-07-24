package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"
)

func (repo *repo) CreateVariant(ctx context.Context, toCreate *models.ProductVariant) (*models.ProductVariant, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.CreateVariant")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).CreateProductVariant(ctx, sqlc.CreateProductVariantParams{
		EditionID:   toCreate.EditionID,
		ProductID:   toCreate.ProductID,
		VendorCode:  toCreate.VendorCode,
		Name:        toCreate.Name,
		Description: toCreate.Description,
		Price:       toCreate.Price,
		Stock:       toCreate.Stock,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapVariant(result)), nil
}
