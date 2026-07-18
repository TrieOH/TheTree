package ports

import (
	"context"
	"encoding/json"
	"payssage/models"
	"time"
)

// PaymentStatus is a normalized status across providers.
type PaymentStatus string

const (
	StatusPending    PaymentStatus = "pending"
	StatusProcessing PaymentStatus = "processing"
	StatusApproved   PaymentStatus = "approved"
	StatusRejected   PaymentStatus = "rejected"
	StatusCancelled  PaymentStatus = "cancelled"
	StatusRefunded   PaymentStatus = "refunded"
)

//type ChargeRequest struct {
//	Intent models.Intent
//	// Amount in the smallest currency unit (cents / centavos).
//	// Each provider impl is responsible for converting if needed.
//	Amount   int64
//	Currency string // ISO 4217: "BRL", "USD"
//
//	PaymentMethod PaymentMethod
//
//	Description       string
//	ExternalReference string // your internal order/cart ID — maps to metadata in Stripe, external_reference in MP
//
//	Payer Payer
//
//	// RedirectURLs are used by MP's hosted checkout.
//	// Stripe ignores these.
//	RedirectURLs *RedirectURLs
//
//	// PaymentMethod hints — each provider maps these to its own enum.
//	// Leave nil to allow all methods.
//	AllowedMethods []PaymentMethod
//
//	MPSellerToken string
//}

type RedirectURLs struct {
	Success string
	Failure string
	Pending string
}

type PaymentMethod string

const (
	MethodCard   PaymentMethod = "card"
	MethodPix    PaymentMethod = "pix"
	MethodBoleto PaymentMethod = "boleto"
)

// CheckoutSession is what InitiateCheckout returns.
// The frontend switches on ProviderName to decide what to render.
type CheckoutSession struct {
	ProviderName      string
	ExternalReference string
	SessionID         string

	// Stripe: frontend uses this to confirm via Stripe Elements.
	ClientSecret string

	// Mercado Pago: frontend redirects to this URL.
	RedirectURL  string
	PreferenceID string

	ExpiresAt *time.Time
}

// PaymentResult is the normalized output for Charge, Refund, and GetStatus.
type PaymentResult struct {
	TransactionID     string
	ExternalReference string
	Status            PaymentStatus
	Amount            int64
	Currency          string
	RefundedAmount    int64 // > 0 on partial or full refunds
	ProviderRaw       any   // escape hatch: the original provider response
}

// RefundRequest supports both full and partial refunds.
type RefundRequest struct {
	TransactionID string
	Amount        int64  // 0 = full refund
	Reason        string // "duplicate", "fraudulent", "requested_by_customer"
}

// WebhookEvent is the normalized inbound event from either provider.
type WebhookEvent struct {
	Provider          string
	EventType         string // "payment.approved", "payment.refunded", etc.
	TransactionID     string
	ExternalReference string
	Status            PaymentStatus
	Raw               []byte // original payload for auditing
}

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
