package intents

import (
	"encoding/json"

	"context"
	"os"
	"time"

	"github.com/MintzyG/fun"

	"payssage/internal/openapi"
	"payssage/models"
)

func (h *Handlers) HardCreateIntent(ctx context.Context, req openapi.HardCreateIntentRequestObject) (openapi.HardCreateIntentResponseObject, error) {
	if os.Getenv("TEST_MODE") != "true" {
		return nil, fun.Err("test mode only").ServiceUnavailable()
	}
	providerData, err := json.Marshal(req.Body.ProviderData)
	if err != nil {
		return nil, fun.ErrBadRequest("invalid provider_data")
	}
	var metadata json.RawMessage
	if req.Body.Metadata != nil {
		metadata, err = json.Marshal(req.Body.Metadata)
		if err != nil {
			return nil, fun.ErrBadRequest("invalid metadata")
		}
	}
	sandbox := false
	if req.Body.Sandbox != nil {
		sandbox = *req.Body.Sandbox
	}
	intent, err := h.ops.HardCreate(ctx, models.HardCreateIntentInput{
		WalletID:      req.Body.WalletId,
		SellerID:      req.Body.SellerId,
		CollectorID:   req.Body.CollectorId,
		AmountCents:   req.Body.AmountCents,
		Currency:      req.Body.Currency,
		Sandbox:       sandbox,
		Provider:      req.Body.Provider,
		Status:        req.Body.Status,
		ProviderData:  providerData,
		Metadata:      &metadata,
		ExternalID:    req.Body.ExternalId,
		ExternalGroup: req.Body.ExternalGroup,
	})
	if err != nil {
		return nil, err
	}
	return openapi.HardCreateIntent201JSONResponse{
		Code: 201, Data: intent, Timestamp: time.Now(), Module: module,
	}, nil
}
