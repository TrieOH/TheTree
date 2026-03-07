package queries

import (
	"context"
	"univents/internal/commerce/domain"
)

func (uc *QueryService) GetPurchaseByPaymentID(ctx context.Context, paymentID string) (out *domain.Purchase, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "ProductService.GetPurchaseByPaymentID")
	defer span.End()

	return uc.purchases.GetByPaymentID(ctx, paymentID)
}
