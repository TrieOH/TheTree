package errs

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Incoming provider webhooks.
var (
	ErrUnknownWebhookProvider   = errors.New("unknown webhook provider")
	ErrWebhookNotImplemented    = errors.New("webhook handling not implemented for this provider")
	ErrWebhookPayloadInvalid    = errors.New("mercadopago webhook: invalid payload")
	ErrWebhookMissingDataID     = errors.New("mercadopago webhook: missing data.id")
	ErrWebhookSignatureMissing  = errors.New("mercadopago webhook: missing ts or v1 in x-signature header")
	ErrWebhookSignatureMismatch = errors.New("mercadopago webhook: signature mismatch")
)

// MissingParamError is a required header or query parameter the request
// lacks. In is "header" or "query".
type MissingParamError struct {
	In   string
	Name string
}

func (e *MissingParamError) Error() string { return "missing required " + e.In + ": " + e.Name }

// UnhandledWebhookEventError is an event type Payssage does not process;
// the webhook is acknowledged and ignored.
type UnhandledWebhookEventError struct{ Type string }

func (e *UnhandledWebhookEventError) Error() string {
	return fmt.Sprintf("mercadopago webhook: unhandled event type %q", e.Type)
}

// WebhookIntentNotFoundError is a payment no intent was created for; the
// webhook is acknowledged and ignored.
type WebhookIntentNotFoundError struct{ PaymentID string }

func (e *WebhookIntentNotFoundError) Error() string {
	return "mercadopago webhook: no intent found for payment " + e.PaymentID
}

// WebhookSellerTokenMissingError is a webhook for an intent whose seller has
// no access token to fetch the payment with.
type WebhookSellerTokenMissingError struct {
	SellerID uuid.UUID
	IntentID uuid.UUID
}

func (e *WebhookSellerTokenMissingError) Error() string {
	return fmt.Sprintf("seller %s has no mercadopago access token; cannot resolve webhook for intent %s", e.SellerID, e.IntentID)
}
