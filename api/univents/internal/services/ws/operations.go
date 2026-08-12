// Package ws is the per-purchase WebSocket surface (split 6, D9): the only
// live-update channel for a buyer. One socket per purchase — it opens at
// checkout (the ws_token comes in the checkout response) and a disconnect is
// repaired by reconnecting to the same purchase with a fresh token (canonical
// case: a mobile user hops to the bank app to pay pix, comes back). The
// socket delivers a snapshot frame on connect (so a reconnect restores
// context), then live frames until a terminal event (confirmed/expired/
// cancelled) closes it.
//
// The socket is fed by the notifier (LISTEN/NOTIFY, kind="purchase" payloads
// published by the webhook receiver / checkout / expiry) — multi-instance
// safe: each instance's hub subscribes and fans out to its own conns.
package ws

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"univents/internal/services/notify"
	"univents/models"
	"univents/ports"

	payssage "sdk/payssage"

	"github.com/google/uuid"
)

// tokenTTL is how long a handshake token lives. Short-lived because it rides
// the query string (logs): long enough for a checkout → socket open or a
// reconnect hop to the bank app and back.
const tokenTTL = 10 * time.Minute

// tokenBytes is the entropy of the raw token (32 bytes → 43 chars base64url).
const tokenBytes = 32

// Notifier is the LISTEN/NOTIFY subscribe surface (split 6 subscribes to
// kind="purchase"). Satisfied by *database.Notifier (lib/go/database).
type Notifier interface {
	Subscribe(channel string, handler func(payload string))
}

// IntentClient is the Payssage intent read seam, used for the snapshot's
// best-effort intent_status and the cancelled frame's status_detail.
// Satisfied by *payssage.Client and faked in tests.
type IntentClient interface {
	GetIntent(ctx context.Context, intentID uuid.UUID) (*payssage.Intent, error)
}

type Operations struct {
	tokens    ports.WsTokenRepo
	purchases ports.PurchaseRepo
	intents   IntentClient
	hub       *Hub

	// intentAttempts / intentRetryDelay are the snapshot's intent-status
	// fetch budget (mirrors the split-5 resume): up to 3 tries with a short
	// pause, then degrade gracefully — a Payssage blip must not fail the
	// snapshot, the purchases table is the source of truth.
	intentAttempts   int
	intentRetryDelay time.Duration
}

func NewOperations(tokens ports.WsTokenRepo, purchases ports.PurchaseRepo, intents IntentClient, notifier Notifier) *Operations {
	o := &Operations{
		tokens:           tokens,
		purchases:        purchases,
		intents:          intents,
		intentAttempts:   3,
		intentRetryDelay: 250 * time.Millisecond,
	}
	o.hub = newHub(purchases, intents)
	// One subscription for the life of the process; the handler spawns a
	// goroutine so the notifier's listen loop stays free (the fan-out may
	// hit the DB and Payssage).
	notifier.Subscribe(notify.Channel, func(payload string) {
		go o.hub.handleNotification(payload)
	})
	return o
}

// SetIntentRetry overrides the snapshot intent-status retry budget. Test-only.
func (o *Operations) SetIntentRetry(attempts int, delay time.Duration) {
	o.intentAttempts = attempts
	o.intentRetryDelay = delay
}

// IssueToken is getWsToken: the caller must be the purchase's purchaser
// (anything else is NOT_FOUND via the repo — no existence leak). Issues a
// fresh one-time token — a random 32-byte value stored as its SHA-256 hash,
// valid 10 minutes — and returns the raw value (the only time it exists).
// Transaction-aware through the shared tx runner: checkout (split 7) calls
// the same method inside its tx and the token + purchase commit atomically.
func (o *Operations) IssueToken(ctx context.Context, purchaseID, userID uuid.UUID) (string, time.Time, error) {
	_, err := o.purchases.GetByIDForOwner(ctx, purchaseID, userID)
	if err != nil {
		return "", time.Time{}, err
	}

	raw, err := newToken()
	if err != nil {
		return "", time.Time{}, err
	}
	expires := time.Now().UTC().Add(tokenTTL)
	_, err = o.tokens.Create(ctx, &models.WsToken{
		PurchaseID: purchaseID,
		UserID:     userID,
		TokenHash:  hashToken(raw),
		ExpiresAt:  expires,
	})
	if err != nil {
		return "", time.Time{}, err
	}
	return raw, expires, nil
}

// newToken returns a random URL-safe token. crypto/rand failure is a
// genuine entropy failure — surfaced as an error, never a weak token.
func newToken() (string, error) {
	b := make([]byte, tokenBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// hashToken is the at-rest form: SHA-256 hex. The raw token is never
// persisted or logged.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
