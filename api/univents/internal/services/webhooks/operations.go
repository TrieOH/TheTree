// Package webhooks is the Payssage webhook receiver — the only component
// that confirms payment (D3). It verifies the delivery signature, correlates
// the intent to a purchase (D2), and flips the purchase + materialized rows
// on approval, failure, or late approval (D4). Checkout never self-approves.
package webhooks

import (
	"context"
	"time"

	"lib/database"
	"univents/internal/services/notify"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
	"github.com/riverqueue/river/rivertype"
)

// Badges is the badge-emission surface the receiver calls on approval.
// Satisfied by *badges.Operations.
type Badges interface {
	EmitForConfirmedRegistration(ctx context.Context, registrationID uuid.UUID) (*models.BadgeEmission, error)
}

// Notifier is the LISTEN/NOTIFY publish surface (split 6 subscribes).
// Satisfied by *database.Notifier (lib/go/database).
type Notifier interface {
	Notify(ctx context.Context, channel, payload string) error
}

// River is the expiry-job cancellation surface. Satisfied by
// *river.Client[pgx.Tx]. A no-op until split 7 checkout schedules the
// 10:01 expiry job (seeded split-4 purchases carry a null river_job_id).
type River interface {
	JobCancel(ctx context.Context, jobID int64) (*rivertype.JobRow, error)
}

// channelUniventsChanges is the single NOTIFY channel the store publishes
// on (contract in internal/services/notify): the WS hub routes
// kind="purchase" events (D9) and the SSE relay routes kind="stock" deltas
// (D10), both subscribed in split 6.
const channelUniventsChanges = notify.Channel

type Operations struct {
	purchases        ports.PurchaseRepo
	registrations    ports.RegistrationRepo
	productPurchases ports.ProductPurchaseRepo
	participations   ports.ProgramParticipationRepo
	badges           Badges
	notifier         Notifier
	river            River
	tx               database.TxRunner
	secret           string
	// cardRaceWait is the D3 wait before re-querying a purchase whose
	// checkout tx may not have committed yet. Overridable in tests.
	cardRaceWait time.Duration
}

// SetCardRaceWait overrides the D3 wait between the first miss and the
// re-query. Test-only (zero makes the race tests instant).
func (o *Operations) SetCardRaceWait(d time.Duration) { o.cardRaceWait = d }

func NewOperations(
	purchases ports.PurchaseRepo,
	registrations ports.RegistrationRepo,
	productPurchases ports.ProductPurchaseRepo,
	participations ports.ProgramParticipationRepo,
	badges Badges,
	notifier Notifier,
	river River,
	tx database.TxRunner,
	secret string,
) *Operations {
	return &Operations{
		purchases:        purchases,
		registrations:    registrations,
		productPurchases: productPurchases,
		participations:   participations,
		badges:           badges,
		notifier:         notifier,
		river:            river,
		tx:               tx,
		secret:           secret,
		cardRaceWait:     time.Second,
	}
}

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
