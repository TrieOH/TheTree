package ports

import (
	"context"
	"encoding/json"
	"payssage/models"
)

// PaymentAbstractionLayer is the single contract every provider must fulfill.
//
// All methods take *models.Intent and mutate it directly (updating
// ProviderData and any other relevant fields) rather than returning new
// state. Callers are responsible for persisting intent after a successful
// call. This keeps the interface stable across providers that return very
// different shapes of data for the same logical operation — the intent
// itself is the one common surface every provider writes back to.
//
// A non-nil error means the intent was NOT mutated (or, if mutated, should
// not be trusted/persisted) — implementations must not partially mutate
// intent before returning an error.
type PaymentAbstractionLayer interface {
	// Checkout initiates a payment for intent using providerCheckoutData
	// (provider-specific fields, e.g. card token or Pix config). On success,
	// intent.ProviderData is updated with the provider's response state.
	Checkout(ctx context.Context, intent *models.Intent, checkoutData json.RawMessage) error

	// CancelPendingPayment cancels a still-pending payment tied to intent.
	// The provider-specific reference (e.g. MercadoPago transaction ID) is
	// read from intent.ProviderData. On success, intent.ProviderData is
	// updated to reflect the canceled state.
	CancelPendingPayment(ctx context.Context, intent *models.Intent) error
}
