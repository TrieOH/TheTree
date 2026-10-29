package handlers

import (
	"context"

	idx "sdk/identityx"

	"univents/internal/features/checkouts/service"
	"univents/internal/openapi"
	"univents/models"
)

// CreateEditionCheckout is the money path (split 7): reserve the cart and
// create the Payssage intent in one request. The response mirrors the
// shared Purchase read shape (resume == snapshot == checkout uniform) plus
// the one-time ws_token so the front can open the socket immediately.
func (h *Handlers) CreateEditionCheckout(ctx context.Context, req openapi.CreateEditionCheckoutRequestObject) (openapi.CreateEditionCheckoutResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	res, err := h.ops.Checkout(ctx, req.EditionId, ident.Sub.ID, toCheckoutInput(req.Body))
	if err != nil {
		return nil, err
	}

	return openapi.CreateEditionCheckout201JSONResponse(openapi.CheckoutResult{
		PurchaseId:       res.PurchaseID,
		EditionId:        res.EditionID,
		Status:           openapi.PurchaseStatus(res.Status),
		StatusReason:     res.StatusReason,
		TotalCents:       res.TotalCents,
		Currency:         res.Currency,
		PaymentMethod:    res.PaymentMethod,
		PayssageIntentId: res.PayssageIntentID,
		QrCode:           res.QRCode,
		QrCodeBase64:     res.QRCodeBase64,
		ExpiresAt:        res.ExpiresAt,
		CreatedAt:        &res.CreatedAt,
		Items:            toItems(res.Items),
		WsToken:          res.WsToken,
		WsTokenExpiresAt: &res.WsTokenExpiresAt,
	}), nil
}

// toCheckoutInput maps the generated request body to the service's domain
// input. Validation (ticket attendee, per-type quantity, duplicates,
// payment fields for paid orders) runs in the service — the spec only
// enforces the static request shape.
func toCheckoutInput(body *openapi.CreateCheckoutRequest) service.CheckoutInput {
	paymentMethod := ""
	if body.PaymentMethod != nil {
		paymentMethod = string(*body.PaymentMethod)
	}
	payer := service.Payer{}
	if body.Payer != nil {
		payer = service.Payer{
			Email:                body.Payer.Email,
			IdentificationType:   body.Payer.IdentificationType,
			IdentificationNumber: body.Payer.IdentificationNumber,
		}
	}
	in := service.CheckoutInput{
		PaymentMethod:   paymentMethod,
		CardToken:       body.CardToken,
		PaymentMethodID: body.PaymentMethodId,
		Installments:    body.Installments,
		IssuerID:        body.IssuerId,
		Payer:           payer,
		Items:           make([]service.CheckoutLine, 0, len(body.Items)),
	}
	for _, item := range body.Items {
		line := service.CheckoutLine{
			ItemType: models.PurchaseItemType(item.ItemType),
			ItemID:   item.ItemId,
			Quantity: item.Quantity,
		}
		if item.Attendee != nil {
			line.Attendee = &service.Attendee{
				UserID: item.Attendee.UserId,
				Email:  item.Attendee.Email,
				Name:   item.Attendee.Name,
			}
		}
		in.Items = append(in.Items, line)
	}
	return in
}
