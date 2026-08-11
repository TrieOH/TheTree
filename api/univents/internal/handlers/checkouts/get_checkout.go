package checkouts

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
	"univents/models"
)

// GetCheckout is the resume read: the buyer's source of truth for one
// purchase (owner-scoped — a non-owner gets 404 via the service). The
// response mirrors the split-6 WS snapshot shape so the front treats
// resume == snapshot == checkout uniformly.
func (h *Handlers) GetCheckout(ctx context.Context, req openapi.GetCheckoutRequestObject) (openapi.GetCheckoutResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	res, err := h.ops.Get(ctx, req.PurchaseId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.GetCheckout200JSONResponse{
		Code: 200,
		Data: &openapi.Checkout{
			PurchaseId:       res.Purchase.ID,
			EditionId:        res.Purchase.EditionID,
			Status:           openapi.PurchaseStatus(res.Purchase.Status),
			StatusReason:     res.Purchase.StatusReason,
			TotalCents:       res.Purchase.TotalCents,
			Currency:         res.Purchase.Currency,
			PaymentMethod:    res.Purchase.PaymentMethod,
			PayssageIntentId: res.Purchase.PayssageIntentID,
			QrCode:           res.Purchase.QRCode,
			QrCodeBase64:     res.Purchase.QRCodeBase64,
			ExpiresAt:        res.Purchase.ExpiresAt,
			CreatedAt:        &res.Purchase.CreatedAt,
			Items:            toItems(res.Items),
			IntentStatus:     res.IntentStatus,
		},
		Timestamp: time.Now(),
		Module:    module,
	}, nil
}

// toItems maps the ledger rows to the wire shape (item identity + server-
// computed price; materialization links stay internal).
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
