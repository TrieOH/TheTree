package checkouts

import (
	"context"
	"time"

	idx "sdk/identityx"

	"univents/internal/openapi"
)

// ListEditionPurchases is the organizer orders read (owner/admin): every
// purchase of an edition with items + attendees.
func (h *Handlers) ListEditionPurchases(ctx context.Context, req openapi.ListEditionPurchasesRequestObject) (openapi.ListEditionPurchasesResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := h.ops.ListEditionPurchases(ctx, ident.Sub.ID, req.EditionId)
	if err != nil {
		return nil, err
	}
	data := make([]openapi.EditionPurchase, 0, len(rows))
	for _, row := range rows {
		attendees := make([]openapi.PurchaseAttendee, 0, len(row.Attendees))
		for _, a := range row.Attendees {
			attendees = append(attendees, openapi.PurchaseAttendee{Name: a.Name, Email: a.Email})
		}
		data = append(data, openapi.EditionPurchase{
			PurchaseId:       row.Purchase.ID,
			EditionId:        row.Purchase.EditionID,
			Status:           openapi.PurchaseStatus(row.Purchase.Status),
			StatusReason:     row.Purchase.StatusReason,
			TotalCents:       row.Purchase.TotalCents,
			Currency:         row.Purchase.Currency,
			PaymentMethod:    row.Purchase.PaymentMethod,
			PayerEmail:       row.Purchase.PayerEmail,
			PayssageIntentId: row.Purchase.PayssageIntentID,
			CreatedAt:        &row.Purchase.CreatedAt,
			Items:            toItems(row.Items),
			Attendees:        attendees,
		})
	}
	return openapi.ListEditionPurchases200JSONResponse{
		Code: 200, Data: &data, Timestamp: time.Now(), Module: module,
	}, nil
}

// RefundPurchase initiates a full refund of an approved purchase
// (owner/admin). The purchase stays approved until the webhook flips it.
func (h *Handlers) RefundPurchase(ctx context.Context, req openapi.RefundPurchaseRequestObject) (openapi.RefundPurchaseResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	purchase, err := h.ops.RefundPurchase(ctx, ident.Sub.ID, req.PurchaseId)
	if err != nil {
		return nil, err
	}
	return openapi.RefundPurchase200JSONResponse{
		Code: 200,
		Data: &openapi.Purchase{
			PurchaseId:       purchase.ID,
			EditionId:        purchase.EditionID,
			Status:           openapi.PurchaseStatus(purchase.Status),
			StatusReason:     purchase.StatusReason,
			TotalCents:       purchase.TotalCents,
			Currency:         purchase.Currency,
			PaymentMethod:    purchase.PaymentMethod,
			PayssageIntentId: purchase.PayssageIntentID,
			PayerEmail:       purchase.PayerEmail,
			QrCode:           purchase.QRCode,
			QrCodeBase64:     purchase.QRCodeBase64,
			ExpiresAt:        purchase.ExpiresAt,
			CreatedAt:        &purchase.CreatedAt,
			Items:            []openapi.PurchaseItem{},
		},
		Timestamp: time.Now(),
		Module:    module,
	}, nil
}
