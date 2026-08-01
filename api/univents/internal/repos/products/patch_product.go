package products

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) PatchProduct(ctx context.Context, id uuid.UUID, product *models.Product) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsRepo.PatchProduct")
	defer span.End()
	result, err := database.Queries(ctx, repo.q).PatchProduct(ctx, sqlc.PatchProductParams{
		VendorCode:           product.VendorCode,
		RequiresRegistration: product.RequiresRegistration,
		ID:                   id,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProduct(result)), nil
}
