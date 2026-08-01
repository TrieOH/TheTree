package products

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.GetProductByID")
	defer span.End()
	return o.products.GetProductByID(ctx, id)
}
