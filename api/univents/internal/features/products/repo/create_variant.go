package repo

import (
	"context"
	"lib/database"
	"univents/internal/sqlc"
	"univents/models"
)

func (repo *Repo) CreateVariant(ctx context.Context, toCreate *models.ProductVariant) (*models.ProductVariant, error) {
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
