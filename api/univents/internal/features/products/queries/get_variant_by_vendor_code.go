package queries

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (q *Queries) GetVariantByVendorCode(ctx context.Context, editionID uuid.UUID, vendorCode string) (*models.ProductVariant, error) {
	ctx, span := q.tracer.Start(ctx, "ProductsService.GetVariantByVendorCode")
	defer span.End()
	return q.products.GetVariantByVendorCode(ctx, editionID, vendorCode)
}
