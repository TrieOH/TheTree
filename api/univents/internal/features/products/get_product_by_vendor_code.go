package products

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetProductByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.GetProductByVendorCode")
	defer span.End()
	return o.products.GetProductByVendorCode(ctx, editionID, vendorCode)
}
