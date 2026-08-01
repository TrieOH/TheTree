package products

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]models.ProductVariant, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.ListVariantsByProduct")
	defer span.End()
	return o.products.ListVariantsByProduct(ctx, productID)
}
