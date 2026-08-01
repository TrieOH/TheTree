package intents

import (
	"encoding/json"

	"context"
	"time"

	"github.com/MintzyG/fun"

	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) Checkout(ctx context.Context, req openapi.CheckoutRequestObject) (openapi.CheckoutResponseObject, error) {
	checkoutData, err := json.Marshal(req.Body.CheckoutProviderData)
	if err != nil {
		return nil, fun.ErrBadRequest("invalid checkout_provider_data")
	}
	var metadata json.RawMessage
	if req.Body.Metadata != nil {
		metadata, err = json.Marshal(req.Body.Metadata)
		if err != nil {
			return nil, fun.ErrBadRequest("invalid metadata")
		}
	}
	intent, err := h.ops.Checkout(ctx, models.CreateIntentInput{
		WalletID:     req.WalletId,
		SellerID:     req.Body.SellerId,
		Currency:     req.Body.Currency,
		AmountCents:  req.Body.AmountCents,
		CheckoutData: checkoutData,
		Metadata:     &metadata,
	})
	if err != nil {
		return nil, err
	}
	return openapi.Checkout201JSONResponse{
		Code: 201, Data: intent, Timestamp: time.Now(), Module: module,
	}, nil
}
