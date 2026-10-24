// Package webhooks is the Payssage webhook receiver — the only component
// that confirms payment (D3). It verifies the delivery signature, correlates
// the intent to a purchase (D2), and flips the purchase + materialized rows
// on approval, failure, or late approval (D4). Checkout never self-approves.
package service

import (
	"context"
	"time"

	"lib/database"
	libjobs "lib/jobs"
	"univents/internal/notify"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
)

// Badges is the badge-emission surface the receiver calls on approval and
// refund. Satisfied by *badges.Operations.
type Badges interface {
	EmitForConfirmedRegistration(ctx context.Context, registrationID uuid.UUID) (*models.BadgeEmission, error)
	// RevokeForRegistration revokes a participant badge when its
	// registration is cancelled (refund plan slice 3: refunds flip
	// registrations to cancelled and revoke the emitted badge).
	RevokeForRegistration(ctx context.Context, registrationID uuid.UUID, reason string) error
}

// Notifier is the LISTEN/NOTIFY publish surface (split 6 subscribes).
// Satisfied by *database.Notifier (lib/go/database).
type Notifier interface {
	Notify(ctx context.Context, channel, payload string) error
}

// Jobs is the background-work surface: the gifted-ticket email
// (gifts.send_email) is enqueued inside the approve tx, before commit, with
// a delay so it cannot run before the confirmation is visible. Approval
// does not cancel the purchase's expiry job — a queued message cannot be
// withdrawn; the expiry worker no-ops on a purchase that is no longer
// pending. Satisfied by the lib/jobs carriers.
type Jobs interface {
	EnqueueAfter(ctx context.Context, job libjobs.Job, delay time.Duration) error
}

// channelUniventsChanges is the single NOTIFY channel the store publishes
// on (contract in internal/notify): the WS hub routes
// kind="purchase" events (D9) and the SSE relay routes kind="stock" deltas
// (D10), both subscribed in split 6.
const channelUniventsChanges = notify.Channel

type Deps struct {
	Purchases        ports.PurchaseRepo
	Registrations    ports.RegistrationRepo
	ProductPurchases ports.ProductPurchaseRepo
	Participations   ports.ProgramParticipationRepo
	Badges           Badges
	Notifier         Notifier
	Jobs             Jobs
	Tx               database.TxRunner
	Secret           string
}

type Operations struct {
	Deps

	// cardRaceWait is the D3 wait before re-querying a purchase whose
	// checkout tx may not have committed yet. Overridable in tests.
	cardRaceWait time.Duration
}

func New(d Deps) *Operations {
	return &Operations{
		Deps:         d,
		cardRaceWait: time.Second,
	}
}

// SetCardRaceWait overrides the D3 wait between the first miss and the
// re-query. Test-only (zero makes the race tests instant).
func (o *Operations) SetCardRaceWait(d time.Duration) { o.cardRaceWait = d }

// ReceiveInput is one Payssage webhook delivery, reduced to what the
// receiver acts on: correlation key (intent_id), event type, and the raw
// body + signature header for verification.
type ReceiveInput struct {
	IntentID     uuid.UUID
	EventType    string
	RawBody      []byte
	Signature    string
	StatusReason *string // normalized outcome detail from the payssage envelope (e.g. "high_risk")
}
