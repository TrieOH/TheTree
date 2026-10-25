package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetVariantByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.ProductVariant, error) {
	return o.Products.GetVariantByVendorCode(ctx, editionID, vendorCode)
}
