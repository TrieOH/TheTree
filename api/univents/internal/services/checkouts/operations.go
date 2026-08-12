// Package checkouts serves the edition checkout (split 7 — the money path)
// and the resume read (split 5). Checkout reserves the cart and creates the
// Payssage intent in one request; the webhook receiver (split 4) and the
// expiry worker (split 7) own all status changes — checkout never
// self-approves (D3).
package checkouts

import (
	"context"
	"time"

	"lib/database"
	"univents/internal/services/notify"
	"univents/models"
	"univents/ports"

	payssage "sdk/payssage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// defaultIntentAttempts / defaultIntentRetryDelay are the resume's intent
// fetch budget: up to 3 tries with a short pause, then degrade gracefully
// (a Payssage blip must not fail the source-of-truth read).
const (
	defaultIntentAttempts   = 3
	defaultIntentRetryDelay = 250 * time.Millisecond
)

// IntentClient is the Payssage intent read seam (resume's `intent_status`).
// Satisfied by *payssage.Client (sdk/go/Payssage) and faked in tests.
type IntentClient interface {
	GetIntent(ctx context.Context, intentID uuid.UUID) (*payssage.Intent, error)
}

// PayssageClient is the checkout's Payssage write seam: create the intent
// (post-commit) and best-effort cancel it when storing the intent id back
// fails. Satisfied by *payssage.Client and faked in tests.
type PayssageClient interface {
	Checkout(ctx context.Context, walletID uuid.UUID, req payssage.CreateIntentRequest) (*payssage.Intent, error)
	CancelIntent(ctx context.Context, intentID uuid.UUID) (*payssage.Intent, error)
}

// TokenIssuer is the one-time WS handshake token seam (split 6). Satisfied
// by *ws.Operations; transaction-aware through the shared tx runner, so the
// token + purchase commit atomically inside the checkout tx.
type TokenIssuer interface {
	IssueToken(ctx context.Context, purchaseID, userID uuid.UUID) (string, time.Time, error)
}

// Badges is the badge-emission surface the free-order path calls (a free
// ticket is confirmed at checkout, so the participant badge emits
// immediately). Satisfied by *badges.Operations.
type Badges interface {
	EmitForConfirmedRegistration(ctx context.Context, registrationID uuid.UUID) (*models.BadgeEmission, error)
}

// Notifier is the LISTEN/NOTIFY publish surface (split 6 subscribes).
// Satisfied by *database.Notifier (lib/go/database).
type Notifier interface {
	Notify(ctx context.Context, channel, payload string) error
}

// River is the expiry-job scheduling surface: river.InsertTx inside the
// checkout tx (payssage's dispatchDeliveries pattern) so the job + purchase
// commit atomically — no orphan job on rollback, no job before the purchase
// is visible. Satisfied by *river.Client[pgx.Tx].
type River interface {
	InsertTx(ctx context.Context, tx pgx.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// channelUniventsChanges is the single NOTIFY channel the store publishes
// on (contract in internal/services/notify).
const channelUniventsChanges = notify.Channel

type Operations struct {
	purchases        ports.PurchaseRepo
	editions         ports.EditionRepo
	events           ports.EventRepo
	ticketTypes      ports.TicketTypeRepo
	products         ports.ProductRepo // product variants (variants live on the products repo)
	programs         ports.ProgramRepo
	occurrences      ports.ProgramOccurrenceRepo
	registrations    ports.RegistrationRepo
	productPurchases ports.ProductPurchaseRepo
	participations   ports.ProgramParticipationRepo
	badges           Badges
	notifier         Notifier
	river            River
	tx               database.TxRunner
	intents          IntentClient // resume (split 5)
	payssage         PayssageClient
	walletID         uuid.UUID // the single platform wallet (D6)
	tokens           TokenIssuer

	// intentAttempts / intentRetryDelay are the resume's GetIntent budget.
	// Overridable in tests.
	intentAttempts   int
	intentRetryDelay time.Duration
}

func NewOperations(
	purchases ports.PurchaseRepo,
	editions ports.EditionRepo,
	events ports.EventRepo,
	ticketTypes ports.TicketTypeRepo,
	products ports.ProductRepo,
	programs ports.ProgramRepo,
	occurrences ports.ProgramOccurrenceRepo,
	registrations ports.RegistrationRepo,
	productPurchases ports.ProductPurchaseRepo,
	participations ports.ProgramParticipationRepo,
	badges Badges,
	notifier Notifier,
	river River,
	tx database.TxRunner,
	intents IntentClient,
	payssage PayssageClient,
	walletID uuid.UUID,
	tokens TokenIssuer,
) *Operations {
	return &Operations{
		purchases:        purchases,
		editions:         editions,
		events:           events,
		ticketTypes:      ticketTypes,
		products:         products,
		programs:         programs,
		occurrences:      occurrences,
		registrations:    registrations,
		productPurchases: productPurchases,
		participations:   participations,
		badges:           badges,
		notifier:         notifier,
		river:            river,
		tx:               tx,
		intents:          intents,
		payssage:         payssage,
		walletID:         walletID,
		tokens:           tokens,
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
