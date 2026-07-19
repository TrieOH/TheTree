package ports

import (
	"context"
	"net/http"
	"payssage/models"
)

// WebhookProvider is the per-provider contract for verifying and parsing
// an inbound webhook. Each provider verifies authenticity in whatever way
// it requires (HMAC over a manifest, raw-body HMAC, RSA signature, etc.).
// Implementations read whatever they need off r (headers, query, body).
type WebhookProvider interface {
	// VerifySignature checks the raw inbound request against the
	// provider's own signing scheme, using credentials/secrets it manages
	// internally (env config, not passed in — these are integrator-level,
	// not per-wallet). Must be called before Parse.
	VerifySignature(ctx context.Context, r *http.Request, rawBody []byte) error

	// Parse extracts a normalized event from an already-verified webhook
	// request. Providers with "thin" payloads (MercadoPago sends only a
	// type + resource id) are expected to call back into the provider's
	// API inside Parse to resolve full state.
	Parse(ctx context.Context, r *http.Request, rawBody []byte) (*models.WebhookParseResult, error)
}
