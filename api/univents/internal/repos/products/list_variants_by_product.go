package products

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *Repo) ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]models.ProductVariant, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsRepo.ListVariantsByProduct")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListProductVariantsByProduct(ctx, productID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapVariant), nil
}
