// Package checkouts serves the resume read (split 5): one purchase with its
// items and — when the purchase is still pending and carries a Payssage
// intent — the intent's current status, so a resume surfaces
// `intent.updated` without polling. Strictly read-only: the webhook
// receiver (split 4) and the expiry worker (split 7) own status changes.
package checkouts

import (
	"context"
	"time"
	"univents/models"
	"univents/ports"

	"sdk/payssage"

	"github.com/google/uuid"
)

// defaultIntentAttempts / defaultIntentRetryDelay are the resume's intent
// fetch budget: up to 3 tries with a short pause, then degrade gracefully
// (a Payssage blip must not fail the source-of-truth read).
const (
	defaultIntentAttempts   = 3
	defaultIntentRetryDelay = 250 * time.Millisecond
)

// IntentClient is the Payssage intent read seam. Satisfied by
// *payssage.Client (sdk/go/Payssage) and faked in tests.
type IntentClient interface {
	GetIntent(ctx context.Context, intentID uuid.UUID) (*payssage.Intent, error)
}

type Operations struct {
	purchases ports.PurchaseRepo
	intents   IntentClient
	// intentAttempts is how many GetIntent tries the resume makes before
	// degrading; intentRetryDelay is the pause between tries. Overridable
	// in tests.
	intentAttempts   int
	intentRetryDelay time.Duration
}

func NewOperations(purchases ports.PurchaseRepo, intents IntentClient) *Operations {
	return &Operations{
		purchases:        purchases,
		intents:          intents,
		intentAttempts:   defaultIntentAttempts,
		intentRetryDelay: defaultIntentRetryDelay,
	}
}

// SetIntentRetry overrides the intent-status retry budget. Test-only.
func (o *Operations) SetIntentRetry(attempts int, delay time.Duration) {
	o.intentAttempts = attempts
	o.intentRetryDelay = delay
}

// Resume is the getCheckout payload: the purchase (flat, mirroring the WS
// snapshot shape) plus its items and the intent's status when relevant.
type Resume struct {
	Purchase     models.Purchase
	Items        []models.PurchaseItem
	IntentStatus *string
}
