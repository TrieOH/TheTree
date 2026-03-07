package queries

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"

	"github.com/google/uuid"
)

func (uc *QueryService) ListPurchaseItems(ctx context.Context, purchaseID uuid.UUID) (out []domain.LineItem, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "PurchaseService.ListPurchaseItems")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	return uc.purchases.ListPurchaseItems(ctx, purchaseID, sub.ID)
}
