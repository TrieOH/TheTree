package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"univents/models"

	"github.com/google/uuid"
)

func (repo *repo) ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]models.ProductVariant, error) {
	ctx, span := repo.tracer.Start(ctx, "ProductsRepo.ListVariantsByProduct")
	defer span.End()
	results, err := database.Queries(ctx, repo.q).ListProductVariantsByProduct(ctx, productID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(results, mapVariant), nil
}
