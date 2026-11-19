// Package checkouts is the checkouts feature: its handlers and operations.
// Build wires it.
//
// The operations serve the edition checkout (split 7 — the money path)
// and the resume read (split 5). Checkout reserves the cart and creates the
// Payssage intent in one request; the webhook receiver (split 4) and the
// expiry worker (split 7) own all status changes — checkout never
// self-approves (D3).
package checkouts

import (
	"context"
	"errors"
	"lib/database"
	libjobs "lib/jobs"
	"lib/telemetry"
	"log/slog"
	idx "sdk/identityx"
	payssage "sdk/payssage"
	"strings"
	"time"
	"univents/errs"
	"univents/internal/authz"
	"univents/internal/notify"
	"univents/models"
	"univents/ports"
	"uuid"
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
// (post-commit), best-effort cancel it when storing the intent id back
// fails, and refund it on the organizer refund path. Satisfied by
// *payssage.Client and faked in tests.
type PayssageClient interface {
	Checkout(ctx context.Context, walletID uuid.UUID, req payssage.CreateIntentRequest) (*payssage.Intent, error)
	CancelIntent(ctx context.Context, intentID uuid.UUID) (*payssage.Intent, error)
	RefundIntent(ctx context.Context, intentID uuid.UUID) (*payssage.Intent, error)
}

// TokenIssuer is the one-time WS handshake token seam (split 6). Satisfied
// by *ws.Operations; transaction-aware through the shared tx runner, so the
// token + purchase commit atomically inside the checkout tx.
type TokenIssuer interface {
	IssueToken(ctx context.Context, purchaseID, userID uuid.UUID) (string, time.Time, error)
}

// ActorResolver is the IdentityX actor lookup seam (ports.ActorResolver),
// used by checkout for attendee resolution (gift buying) and the gift claim.
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

// Jobs is the background-work surface: the expiry and gift-email jobs are
// enqueued inside the checkout tx, before commit, with a delay — the job
// cannot run before the purchase is visible, a failed enqueue fails the
// checkout, and a rolled-back checkout leaves a job that finds nothing and
// no-ops. Satisfied by the lib/jobs carriers.
type Jobs interface {
	EnqueueAfter(ctx context.Context, job libjobs.Job, delay time.Duration) error
}

// channelUniventsChanges is the single NOTIFY channel the store publishes
// on (contract in internal/notify).
const channelUniventsChanges = notify.Channel

type Deps struct {
	Purchases        ports.PurchaseRepo
	Editions         ports.EditionRepo
	Events           ports.EventRepo
	TicketTypes      ports.TicketTypeRepo
	Products         ports.ProductRepo // product variants (variants live on the products repo)
	Programs         ports.ProgramRepo
	Occurrences      ports.ProgramOccurrenceRepo
	Registrations    ports.RegistrationRepo
	ProductPurchases ports.ProductPurchaseRepo
	Participations   ports.ProgramParticipationRepo
	Badges           Badges
	Notifier         Notifier
	Jobs             Jobs
	Tx               database.TxRunner
	Intents          IntentClient // resume (split 5)
	Payssage         PayssageClient
	WalletID         uuid.UUID // the single platform wallet (D6)
	Tokens           TokenIssuer
	Actors           ports.ActorResolver // attendee email → account (gift buying)
	Authz            *authz.Service      // organizer orders/refund (refund plan B3)
}

type Operations struct {
	Deps

	// intentAttempts / intentRetryDelay are the resume's GetIntent budget.
	// Overridable in tests.
	intentAttempts   int
	intentRetryDelay time.Duration
}

func New(d Deps) *Operations {
	return &Operations{
		Deps:             d,
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

// sdkActors adapts the IdentityX SDK ActorService to the ActorResolver
// seam (ports.ActorResolver): the SDK surfaces "no account for this email"
// / "unknown actor" as a 404 APIError, which is mapped to
// ErrActorNotFound so the checkout logic can branch on "no account yet"
// (email-only gift) vs "verification failed" (error).
type sdkActors struct {
	actors *idx.ActorService
}

// NewSDKActorResolver wraps the IdentityX actor service for checkout's
// attendee resolution, gift claim, and the profile-badges gift claim.
// Lookups are scoped to the univents project (the client's project id).
func NewSDKActorResolver(actors *idx.ActorService) ports.ActorResolver {
	return &sdkActors{actors: actors}
}

func (s *sdkActors) GetByEmail(ctx context.Context, email string) (*idx.Actor, error) {
	actor, err := s.actors.GetByEmail(ctx, email)
	if err != nil {
		if idx.IsNotFound(err) {
			return nil, ports.ErrActorNotFound
		}
		return nil, err
	}
	return actor, nil
}

func (s *sdkActors) GetByID(ctx context.Context, id uuid.UUID) (*idx.Actor, error) {
	actor, err := s.actors.GetByID(ctx, id)
	if err != nil {
		if idx.IsNotFound(err) {
			return nil, ports.ErrActorNotFound
		}
		return nil, err
	}
	return actor, nil
}

// AttendeeCount returns the number of confirmed attendees of an edition —
// the public attendee-count read (GET /editions/{id}/attendees/count), the
// storefront's "N people already registered" number. Only paid
// registrations count (status confirmed); pending reservations and
// cancelled/expired rows do not. Unknown editions are NOT_FOUND (the
// endpoint exists, the edition does not).
func (o *Operations) AttendeeCount(ctx context.Context, editionID uuid.UUID) (int64, error) {
	_, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return 0, err
	}
	return o.Registrations.CountConfirmedByEdition(ctx, editionID)
}

// Get is the resume read: the purchase owner-scoped (the caller must be the
// purchase's purchaser_id — anything else is NOT_FOUND via the repo, no
// existence leak), its items from the availability ledger, and — only when
// the purchase is still pending and carries an intent — the intent's
// current status. Never transitions state.
func (o *Operations) Get(ctx context.Context, purchaseID, purchaserID uuid.UUID) (*Resume, error) {
	purchase, err := o.Purchases.GetByIDForOwner(ctx, purchaseID, purchaserID)
	if err != nil {
		return nil, err
	}
	items, err := o.Purchases.ListItemsByPurchase(ctx, purchaseID)
	if err != nil {
		return nil, err
	}

	resume := &Resume{Purchase: *purchase, Items: items}
	if purchase.Status == models.PurchaseStatusPending && purchase.PayssageIntentID != nil {
		resume.IntentStatus = o.intentStatus(ctx, *purchase.PayssageIntentID)
	}
	return resume, nil
}

// intentStatus fetches the intent's status with the retry budget in
// Operations (up to 3 tries by default), then degrades gracefully: on
// persistent failure it logs loudly — a resume that cannot reach Payssage
// still returns the purchase's own state, because the purchases table is
// the source of truth and `intent_status` is only a supplement.
func (o *Operations) intentStatus(ctx context.Context, intentID uuid.UUID) *string {
	var lastErr error
	for attempt := 1; attempt <= o.intentAttempts; attempt++ {
		intent, err := o.Intents.GetIntent(ctx, intentID)
		if err == nil {
			return new(string(intent.Status))
		}
		lastErr = err
		if attempt < o.intentAttempts {
			select {
			case <-ctx.Done():
				telemetry.Log().Warn("checkouts: intent status fetch cancelled",
					slog.String("intent_id", intentID.String()),
					telemetry.Err(ctx.Err()))
				return nil
			case <-time.After(o.intentRetryDelay):
			}
		}
	}
	telemetry.Log().Error("checkouts: intent status unavailable after retries — degrading (purchases table remains the source of truth)",
		slog.String("intent_id", intentID.String()),
		slog.Int("attempts", o.intentAttempts),
		telemetry.Err(lastErr))
	return nil
}

// MyTicket returns the caller's active ticket for the edition — their own
// registration where attendee_user_id = caller and status is pending or
// confirmed — with its ticket type. Returns (nil, nil) when the caller
// holds none (cancelled/expired registrations do not count). Read-only;
// never transitions state. An unknown edition is NOT_FOUND (the endpoint
// exists, the edition does not).
//
// This is also the lazy gift-claim surface: a caller who holds no ticket
// under their id may be a gifted recipient who just created their account
// — if their account email matches an email-only gift in this edition, the
// registration is tied to their id (the claim) and, when confirmed, the
// deferred badge is emitted (the gift email told them to create the account
// under that email, so the email is the proof).
func (o *Operations) MyTicket(ctx context.Context, editionID, userID uuid.UUID) (*models.MyTicket, error) {
	_, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}

	reg, err := o.activeRegistrationFor(ctx, editionID, userID)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		reg, err = o.claimGiftByAccountEmail(ctx, editionID, userID)
		if err != nil {
			return nil, err
		}
		if reg == nil {
			//nolint:nilnil // no ticket is a normal state (front falls back to the buy flow)
			return nil, nil
		}
	}

	ticketType, err := o.TicketTypes.GetByID(ctx, reg.TicketTypeID)
	if err != nil {
		return nil, err
	}
	return &models.MyTicket{
		RegistrationID: reg.ID,
		TicketType:     *ticketType,
		Status:         reg.Status,
	}, nil
}

// claimGiftByAccountEmail claims the caller's email-only gifted
// registration in the edition, if any: the caller's account email (from
// IdentityX) is matched against active gifts (attendee_user_id NULL) and
// the registration is tied to their actor id. A confirmed claim emits the
// deferred badge. Idempotent: once claimed, the id is set and the my-ticket
// read finds the registration directly (this path no longer runs). No gift
// → nil, nil.
func (o *Operations) claimGiftByAccountEmail(ctx context.Context, editionID, userID uuid.UUID) (*models.Registration, error) {
	actor, err := o.Actors.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ports.ErrActorNotFound) {
			//nolint:nilnil // unknown actor — nothing to claim by
			return nil, nil
		}
		return nil, err
	}
	if actor.Email == nil {
		//nolint:nilnil // no email on the account — cannot match a gift
		return nil, nil
	}
	email := strings.TrimSpace(strings.ToLower(*actor.Email))
	if email == "" {
		//nolint:nilnil // empty email — cannot match a gift
		return nil, nil
	}

	claimed, err := o.Registrations.ClaimByAttendeeEmail(ctx, editionID, email, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			//nolint:nilnil // no active gift for this email — a normal state
			return nil, nil
		}
		return nil, err
	}

	if claimed.Status == models.RegistrationStatusConfirmed {
		// The gift is paid — emit the badge that was deferred at approval
		// (EmitForConfirmedRegistration skips registrations with no
		// account; now that the id is tied, it emits + emails).
		_, err := o.Badges.EmitForConfirmedRegistration(ctx, claimed.ID)
		if err != nil {
			return nil, err
		}
	}
	return claimed, nil
}

// PurchaseAttendee is one person a ticket line was assigned to — surfaced on
// the organizer orders read so the admin can see who (besides the payer) the
// purchase covers. Gifted tickets show the recipient; the payer_email on the
// purchase identifies who gets the refund.
type PurchaseAttendee struct {
	Name  string
	Email string
}

// EditionPurchase is the organizer orders read row: the shared purchase
// shape plus its items and the ticket attendees (from registrations).
type EditionPurchase struct {
	Purchase  models.Purchase
	Items     []models.PurchaseItem
	Attendees []PurchaseAttendee
}

// ListEditionPurchases is the organizer orders read (refund plan B3): every
// purchase of an edition, newest first, with items and attendee names.
// Owner/admin only — non-members 404, staff 403 (via CheckEvent).
func (o *Operations) ListEditionPurchases(ctx context.Context, actorID, editionID uuid.UUID) ([]EditionPurchase, error) {
	edition, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}
	err = o.Authz.EventFor(ctx, actorID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	purchases, err := o.Purchases.ListByEdition(ctx, editionID)
	if err != nil {
		return nil, err
	}
	out := make([]EditionPurchase, 0, len(purchases))
	for _, p := range purchases {
		ep := EditionPurchase{Purchase: p}
		ep.Items, err = o.Purchases.ListItemsByPurchase(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		ep.Attendees, err = o.attendeesFor(ctx, ep.Items)
		if err != nil {
			return nil, err
		}
		out = append(out, ep)
	}
	return out, nil
}

// attendeesFor resolves the ticket items' registrations to attendee
// name/email (deduped — a purchase can have several ticket lines).
func (o *Operations) attendeesFor(ctx context.Context, items []models.PurchaseItem) ([]PurchaseAttendee, error) {
	var out []PurchaseAttendee
	seen := make(map[uuid.UUID]bool)
	for _, item := range items {
		if item.ItemType != models.PurchaseItemTypeTicket || item.RegistrationID == nil {
			continue
		}
		if seen[*item.RegistrationID] {
			continue
		}
		seen[*item.RegistrationID] = true
		reg, err := o.Registrations.GetByID(ctx, *item.RegistrationID)
		if err != nil {
			return nil, err
		}
		out = append(out, PurchaseAttendee{Name: reg.AttendeeName, Email: reg.AttendeeEmail})
	}
	return out, nil
}

// RefundPurchase initiates a full refund of an approved purchase (refund
// plan B3): owner/admin only. It calls payssage RefundIntent and returns the
// purchase — which stays `approved` until the payment.refunded webhook flips
// it (webhook-confirmed, single writer, D-2). An already-refunded purchase
// is a 409; a pending/expired/cancelled one cannot be refunded.
func (o *Operations) RefundPurchase(ctx context.Context, actorID, purchaseID uuid.UUID) (*models.Purchase, error) {
	purchase, err := o.Purchases.GetByID(ctx, purchaseID)
	if err != nil {
		return nil, err
	}
	edition, err := o.Editions.GetByID(ctx, purchase.EditionID)
	if err != nil {
		return nil, err
	}
	err = o.Authz.EventFor(ctx, actorID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	switch purchase.Status {
	case models.PurchaseStatusApproved:
		// ok
	case models.PurchaseStatusRefunded:
		return nil, errs.ErrAlreadyRefunded
	default:
		return nil, &errs.PurchaseStatusError{Status: string(purchase.Status)}
	}
	if purchase.PayssageIntentID == nil {
		return nil, errs.ErrFreeOrderRefund
	}

	_, err = o.Payssage.RefundIntent(ctx, *purchase.PayssageIntentID)
	if err != nil {
		return nil, mapPayssageError(err)
	}
	// The webhook flips the purchase to refunded; return the current state.
	return purchase, nil
}
