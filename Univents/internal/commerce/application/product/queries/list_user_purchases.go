package queries

import (
	"context"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"
)

func (uc *QueryService) ListUserPurchases(ctx context.Context) (out []domain.Purchase, err error) { // FIXME Pagination
	ctx, span := uc.tracer.Start(ctx, "PurchaseService.ListUserPurchases")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	return uc.purchases.ListUserPurchases(ctx, sub.ID)
}
