package purchases

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
	"univents/models"
)

// ListMyPurchases lists the authenticated user's purchases with their
// items, newest first. No pagination yet — the list is unbounded for now
// and pages later (decision taken at review).
func (h *Handlers) ListMyPurchases(ctx context.Context, _ openapi.ListMyPurchasesRequestObject) (openapi.ListMyPurchasesResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	details, err := h.ops.ListForUser(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	list := make([]openapi.Purchase, 0, len(details))
	for _, d := range details {
		list = append(list, toPurchase(d.Purchase, d.Items))
	}
	return openapi.ListMyPurchases200JSONResponse{
		Code:      200,
		Data:      &openapi.MyPurchases{Purchases: list},
		Timestamp: time.Now(),
		Module:    module,
	}, nil
}

// toPurchase maps a purchase + its ledger rows to the shared wire shape.
func toPurchase(p models.Purchase, items []models.PurchaseItem) openapi.Purchase {
	return openapi.Purchase{
		PurchaseId:       p.ID,
		EditionId:        p.EditionID,
		Status:           openapi.PurchaseStatus(p.Status),
		StatusReason:     p.StatusReason,
		TotalCents:       p.TotalCents,
		Currency:         p.Currency,
		PaymentMethod:    p.PaymentMethod,
		PayssageIntentId: p.PayssageIntentID,
		QrCode:           p.QRCode,
		QrCodeBase64:     p.QRCodeBase64,
		ExpiresAt:        p.ExpiresAt,
		CreatedAt:        &p.CreatedAt,
		Items:            toItems(items),
	}
}

func toItems(items []models.PurchaseItem) []openapi.PurchaseItem {
	out := make([]openapi.PurchaseItem, 0, len(items))
	for _, item := range items {
		out = append(out, openapi.PurchaseItem{
			ItemType:       openapi.PurchaseItemType(item.ItemType),
			ItemId:         item.ItemID,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
		})
	}
	return out
}
