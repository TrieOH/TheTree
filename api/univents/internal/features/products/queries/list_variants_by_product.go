package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]models.ProductVariant, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.ListVariantsByProduct")
	defer span.End()
	return q.products.ListVariantsByProduct(ctx, productID)
}
