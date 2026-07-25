package queries

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetProductByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.Product, error) {
	ctx, span := telemetry.StartSpan(ctx, "ProductsService.GetProductByVendorCode")
	defer span.End()
	return q.products.GetProductByVendorCode(ctx, editionID, vendorCode)
}
