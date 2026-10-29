package handlers

import (
	"context"
	"encoding/json"
	"os"

	"payssage/errs"
	"payssage/internal/features/intents/service"
	"payssage/internal/openapi"
	"payssage/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

func (h *Handlers) CancelIntent(ctx context.Context, req openapi.CancelIntentRequestObject) (openapi.CancelIntentResponseObject, error) {
	intent, err := h.ops.Cancel(ctx, req.IntentId)
	if err != nil {
		return nil, err
	}
	return openapi.CancelIntent200JSONResponse(*intent), nil
}

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

func (h *Handlers) GetIntent(ctx context.Context, req openapi.GetIntentRequestObject) (openapi.GetIntentResponseObject, error) {
	intent, err := h.ops.GetByID(ctx, req.IntentId)
	if err != nil {
		return nil, err
	}
	return openapi.GetIntent200JSONResponse(*intent), nil
}

func (h *Handlers) HardCreateIntent(ctx context.Context, req openapi.HardCreateIntentRequestObject) (openapi.HardCreateIntentResponseObject, error) {
	if os.Getenv("TEST_MODE") != "true" {
		return nil, errs.ErrTestModeOnly
	}
	providerData, err := json.Marshal(req.Body.ProviderData)
	if err != nil {
		return nil, errs.ErrInvalidProviderData
	}
	var metadata json.RawMessage
	if req.Body.Metadata != nil {
		metadata, err = json.Marshal(req.Body.Metadata)
		if err != nil {
			return nil, errs.ErrInvalidMetadata
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
	return openapi.HardCreateIntent201JSONResponse(*intent), nil
}

func (h *Handlers) ListIntentsByProfile(ctx context.Context, _ openapi.ListIntentsByProfileRequestObject) (openapi.ListIntentsByProfileResponseObject, error) {
	intents, err := h.ops.ListByProfile(ctx)
	if err != nil {
		return nil, err
	}
	return openapi.ListIntentsByProfile200JSONResponse{Items: intents}, nil
}

func (h *Handlers) ListOrganizationIntents(ctx context.Context, req openapi.ListOrganizationIntentsRequestObject) (openapi.ListOrganizationIntentsResponseObject, error) {
	intents, err := h.ops.ListByOrg(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return openapi.ListOrganizationIntents200JSONResponse{Items: intents}, nil
}

func (h *Handlers) ListWalletIntents(ctx context.Context, req openapi.ListWalletIntentsRequestObject) (openapi.ListWalletIntentsResponseObject, error) {
	intents, err := h.ops.ListByWallet(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	return openapi.ListWalletIntents200JSONResponse{Items: intents}, nil
}

// RefundIntent fully refunds a succeeded intent (webhook confirms the flip).
func (h *Handlers) RefundIntent(ctx context.Context, req openapi.RefundIntentRequestObject) (openapi.RefundIntentResponseObject, error) {
	intent, err := h.ops.Refund(ctx, req.IntentId)
	if err != nil {
		return nil, err
	}
	return openapi.RefundIntent200JSONResponse(*intent), nil
}

// TestmodeRefundIntent simulates a provider refund: flips a succeeded intent
// to refunded without contacting the provider. TEST_MODE-gated like
// hardCreateIntent.
func (h *Handlers) TestmodeRefundIntent(ctx context.Context, req openapi.TestmodeRefundIntentRequestObject) (openapi.TestmodeRefundIntentResponseObject, error) {
	if os.Getenv("TEST_MODE") != "true" {
		return nil, errs.ErrTestModeOnly
	}
	intent, err := h.ops.TestmodeRefund(ctx, req.IntentId)
	if err != nil {
		return nil, err
	}
	return openapi.TestmodeRefundIntent200JSONResponse(*intent), nil
}
