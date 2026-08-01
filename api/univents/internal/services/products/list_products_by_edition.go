package products

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListProductsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.ListProductsByEdition")
	defer span.End()
	return o.products.ListProductsByEdition(ctx, editionID)
}
