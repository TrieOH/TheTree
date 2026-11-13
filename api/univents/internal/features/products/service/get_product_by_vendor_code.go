package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) GetProductByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.Product, error) {
	return o.Products.GetProductByVendorCode(ctx, editionID, vendorCode)
}
