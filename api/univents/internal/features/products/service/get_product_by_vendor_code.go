package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetProductByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.GetProductByVendorCode")
	defer span.End()
	return o.Products.GetProductByVendorCode(ctx, editionID, vendorCode)
}
