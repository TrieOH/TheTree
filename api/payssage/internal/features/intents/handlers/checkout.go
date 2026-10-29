package handlers

import (
	"encoding/json"

	"context"

	"payssage/errs"
	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) Checkout(ctx context.Context, req openapi.CheckoutRequestObject) (openapi.CheckoutResponseObject, error) {
	checkoutData, err := json.Marshal(req.Body.CheckoutProviderData)
	if err != nil {
		return nil, errs.ErrInvalidCheckoutProviderData
	}
	var metadata json.RawMessage
	if req.Body.Metadata != nil {
		metadata, err = json.Marshal(req.Body.Metadata)
		if err != nil {
			return nil, errs.ErrInvalidMetadata
		}
	}
	intent, err := h.ops.Checkout(ctx, models.CreateIntentInput{
		WalletID:      req.WalletId,
		SellerID:      req.Body.SellerId,
		Currency:      req.Body.Currency,
		AmountCents:   req.Body.AmountCents,
		CheckoutData:  checkoutData,
		Metadata:      &metadata,
		ExternalID:    req.Body.ExternalId,
		ExternalGroup: req.Body.ExternalGroup,
	})
	if err != nil {
		return nil, err
	}
	return openapi.Checkout201JSONResponse(*intent), nil
}
