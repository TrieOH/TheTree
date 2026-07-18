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
type PaymentAbstractionLayer interface {
	Checkout(ctx context.Context, intent *models.Intent, checkoutData json.RawMessage) (json.RawMessage, error)

	// Charge performs a direct server-side charge.
	// Use for server-to-server flows where you already have a payment method token.
	Charge(ctx context.Context, intent *models.Intent) (*models.Intent, error)

	// Refund issues a full or partial refund against a prior transaction.
	Refund(ctx context.Context, intent *models.Intent) (*models.Intent, error)
}

// ResolveProvider returns the appropriate PaymentAbstractionLayer for the given request.
// Currently, routes all traffic to MercadoPago as it is the only implemented provider.
//
// When new providers are added, routing logic should consider:
//   - Payment method: Pix above R$80,00 → AbacatePay (flat R$0,80 fee beats MP's 0,99%)
//   - Currency: non-BRL → Stripe
//   - Payment method: card, boleto → MercadoPago
//
// Example future routing:
//
//	if slices.Contains(methods, MethodPix) && amount >= 8000 {
//	    return abacatePayImpl
//	}
//	if currency != "BRL" {
//	    return stripeImpl
//	}
//func ResolveProvider(method PaymentMethod, amount int64, currency string) PaymentAbstractionLayer {
//	return mercadoPagoImpl
//}

type MercadoPagoProvider interface {
	PaymentAbstractionLayer

	// InitiatePixCheckout builds a Pix-specific order payload.
	// Different payment_method structure — no token, no installments.
	// Populates PixQRCode and PixQRCodeB64 on the returned data.
	//InitiatePixCheckout(ctx context.Context, int) (*models.Intent, error)
	//
	//CancelPixCode(ctx context.Context, paymentID string, sellerToken string) error
	//
	//// NormalizeStatus maps MP's order status and status_detail to PaymentStatus.
	//NormalizeStatus(status string, statusDetail string) models.IntentStatus
}
