package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) ListVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]models.ProductVariant, error) {
	return o.Products.ListVariantsByProduct(ctx, productID)
}
