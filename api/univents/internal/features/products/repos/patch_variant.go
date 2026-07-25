package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) PatchVariant(ctx context.Context, id uuid.UUID, variant *models.ProductVariant) (*models.ProductVariant, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsRepo.PatchVariant")
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
