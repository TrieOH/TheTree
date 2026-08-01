// Package intents implements the StrictServerInterface methods for the
// intents feature. The test-mode hard-create is guarded by TEST_MODE.
package intents

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"payssage/internal/openapi"
	"payssage/internal/services"
	"payssage/models"

	"github.com/MintzyG/fun"
)

const module = "Payssage"

type Handlers struct {
	ops *services.Intents
}

func New(ops *services.Intents) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) ListIntentsByProfile(ctx context.Context, _ openapi.ListIntentsByProfileRequestObject) (openapi.ListIntentsByProfileResponseObject, error) {
	intents, err := h.ops.ListByProfile(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListIntentsByProfile200JSONResponse{
		Code: 200, Data: &intents, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) GetIntent(ctx context.Context, req openapi.GetIntentRequestObject) (openapi.GetIntentResponseObject, error) {
	intent, err := h.ops.GetByID(ctx, req.IntentId)
	if err != nil {
		return nil, err
	}
	return openapi.GetIntent200JSONResponse{
		Code: 200, Data: intent, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) CancelIntent(ctx context.Context, req openapi.CancelIntentRequestObject) (openapi.CancelIntentResponseObject, error) {
	intent, err := h.ops.Cancel(ctx, req.IntentId)
	if err != nil {
		return nil, err
	}
	return openapi.CancelIntent200JSONResponse{
		Code: 200, Data: intent, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListWalletIntents(ctx context.Context, req openapi.ListWalletIntentsRequestObject) (openapi.ListWalletIntentsResponseObject, error) {
	intents, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWalletIntents200JSONResponse{
		Code: 200, Data: &intents, Timestamp: time.Now(), Module: module,
	}, nil
}

func (h *Handlers) ListOrganizationIntents(ctx context.Context, req openapi.ListOrganizationIntentsRequestObject) (openapi.ListOrganizationIntentsResponseObject, error) {
	intents, err := h.ops.ListByOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationIntents200JSONResponse{
		Code: 200, Data: &intents, Timestamp: time.Now(), Module: module,
	}, nil
}

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
		WalletID:     req.Body.WalletId,
		SellerID:     req.Body.SellerId,
		CollectorID:  req.Body.CollectorId,
		AmountCents:  req.Body.AmountCents,
		Currency:     req.Body.Currency,
		Sandbox:      sandbox,
		Provider:     req.Body.Provider,
		Status:       req.Body.Status,
		ProviderData: providerData,
		Metadata:     &metadata,
	})
	if err != nil {
		return nil, err
	}
	return openapi.HardCreateIntent201JSONResponse{
		Code: 201, Data: intent, Timestamp: time.Now(), Module: module,
	}, nil
}
