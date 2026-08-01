package products

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetVariantByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.ProductVariant, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.GetVariantByVendorCode")
	defer span.End()
	return o.products.GetVariantByVendorCode(ctx, editionID, vendorCode)
}
