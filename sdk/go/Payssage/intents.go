package payssage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateIntentRequest is the checkout payload: create a payment intent
// against a wallet for the given seller, currency, and amount. Mirrors the
// spec's `CreateIntentRequest` (`checkout_provider_data` is provider-specific
// and may carry payer info, card tokens, etc.). ExternalID/ExternalGroup are
// the caller's correlation ids (e.g. univents purchase id / edition id).
type CreateIntentRequest struct {
	SellerID             uuid.UUID      `json:"seller_id"`
	Currency             string         `json:"currency"`
	AmountCents          int64          `json:"amount_cents"`
	CheckoutProviderData map[string]any `json:"checkout_provider_data,omitempty"`
	Metadata             map[string]any `json:"metadata,omitempty"`
	ExternalID           *string        `json:"external_id,omitempty"`
	ExternalGroup        *string        `json:"external_group,omitempty"`
}

// Checkout starts a checkout: creates a payment intent against the wallet
// for the given seller, currency, and amount. The provider handles the
// payment out-of-band and reports the result via webhook. Pix returns the QR
// in `Intent.ProviderData`; cards charge synchronously (the intent comes
// back `succeeded`).
func (c *Client) Checkout(ctx context.Context, walletID uuid.UUID, req CreateIntentRequest) (*Intent, error) {
	var out Intent
	err := c.do(ctx, "POST", fmt.Sprintf("/wallets/%s/intents", walletID), req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// GetIntent fetches a payment intent by id — the source of truth for a
// checkout's current state (status, status_detail, provider_data).
func (c *Client) GetIntent(ctx context.Context, intentID uuid.UUID) (*Intent, error) {
	var out Intent
	err := c.do(ctx, "GET", fmt.Sprintf("/intents/%s", intentID), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelIntent cancels a pending payment intent and returns its updated
// state. Only intents still `pending` (or otherwise cancellable) can be
// cancelled.
func (c *Client) CancelIntent(ctx context.Context, intentID uuid.UUID) (*Intent, error) {
	var out Intent
	err := c.do(ctx, "POST", fmt.Sprintf("/intents/%s/cancel", intentID), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// RefundIntent fully refunds a succeeded payment intent (full refund only in
// v1 — no amount) and returns its updated state. The intent status stays
// `succeeded` until the payment.refunded webhook confirms.
func (c *Client) RefundIntent(ctx context.Context, intentID uuid.UUID) (*Intent, error) {
	var out Intent
	err := c.do(ctx, "POST", fmt.Sprintf("/intents/%s/refund", intentID), nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
