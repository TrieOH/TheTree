package repos

import (
	"context"
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) PatchVariant(ctx context.Context, id uuid.UUID, variant *models.ProductVariant) (*models.ProductVariant, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.PatchVariant")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).PatchProductVariant(ctx, sqlc.PatchProductVariantParams{
		VendorCode:  variant.VendorCode,
		Name:        variant.Name,
		Description: variant.Description,
		Price:       variant.Price,
		Stock:       variant.Stock,
		ID:          id,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapVariant(result)), nil
}
