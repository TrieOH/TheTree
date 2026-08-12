package checkouts_test

import (
	"context"
	"sync"
	"time"

	"univents/internal/services/checkouts"
	"univents/models"

	payssage "sdk/payssage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// ── External-seam fakes (the DB-backed ports use the real repos) ──────────

// fakePayssage fakes the checkout's Payssage write seam: Checkout (the
// intent creation) and CancelIntent (the Tx-2-failure safety net). Canned
// behavior per test.
type fakePayssage struct {
	mu         sync.Mutex
	checkoutFn func(walletID uuid.UUID, req payssage.CreateIntentRequest) (*payssage.Intent, error)
	cancelled  []uuid.UUID
	intentSeq  int
}

func newFakePayssage(checkoutFn func(uuid.UUID, payssage.CreateIntentRequest) (*payssage.Intent, error)) *fakePayssage {
	return &fakePayssage{checkoutFn: checkoutFn}
}

func (f *fakePayssage) Checkout(_ context.Context, walletID uuid.UUID, req payssage.CreateIntentRequest) (*payssage.Intent, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.intentSeq++
	return f.checkoutFn(walletID, req)
}

func (f *fakePayssage) CancelIntent(_ context.Context, intentID uuid.UUID) (*payssage.Intent, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.cancelled = append(f.cancelled, intentID)
	return &payssage.Intent{ID: intentID, Status: payssage.IntentStatusCancelled}, nil
}

func (f *fakePayssage) cancelCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.cancelled)
}

// fakeBadges records EmitForConfirmedRegistration calls (free-order badge
// emission).
type fakeBadges struct {
	mu      sync.Mutex
	emitted []uuid.UUID
}

func (f *fakeBadges) EmitForConfirmedRegistration(_ context.Context, registrationID uuid.UUID) (*models.BadgeEmission, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.emitted = append(f.emitted, registrationID)
	return &models.BadgeEmission{ID: uuid.New()}, nil
}

func (f *fakeBadges) emittedCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.emitted)
}

// fakeNotifier records every NOTIFY payload (kind=stock / kind=purchase).
type fakeNotifier struct {
	mu       sync.Mutex
	payloads []string
}

func (f *fakeNotifier) Notify(_ context.Context, _ string, payload string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.payloads = append(f.payloads, payload)
	return nil
}

func (f *fakeNotifier) payloadsCopy() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.payloads...)
}

// fakeRiver records river.InsertTx calls (the expiry job scheduling) and
// returns incrementing job ids.
type fakeRiver struct {
	mu       sync.Mutex
	inserted []struct {
		args        river.JobArgs
		scheduledAt time.Time
	}
	nextJobID int64
}

func (f *fakeRiver) InsertTx(_ context.Context, _ pgx.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nextJobID++
	scheduledAt := time.Time{}
	if opts != nil {
		scheduledAt = opts.ScheduledAt
	}
	f.inserted = append(f.inserted, struct {
		args        river.JobArgs
		scheduledAt time.Time
	}{args: args, scheduledAt: scheduledAt})
	return &rivertype.JobInsertResult{Job: &rivertype.JobRow{ID: f.nextJobID}}, nil
}

func (f *fakeRiver) insertCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.inserted)
}

// fakeTokens issues deterministic ws tokens.
type fakeTokens struct {
	mu     sync.Mutex
	issued []uuid.UUID
}

func (f *fakeTokens) IssueToken(_ context.Context, purchaseID, _ uuid.UUID) (string, time.Time, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.issued = append(f.issued, purchaseID)
	return "ws-token-" + purchaseID.String(), time.Now().Add(10 * time.Minute), nil
}

func (f *fakeTokens) issuedCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.issued)
}

// ── Intent builders ────────────────────────────────────────────────────────

// pixIntent returns a pix checkout response with a QR in provider_data.
func pixIntent() *payssage.Intent {
	return &payssage.Intent{
		ID:           uuid.New(),
		Status:       payssage.IntentStatusPending,
		SellerID:     uuid.New(),
		ProviderData: []byte(`{"pix_qr_code":"000201...","pix_qr_code_base64":"iVBORw0KGgo="}`),
	}
}

// cardIntent returns a synchronous card charge: intent already succeeded
// (the purchase stays pending — only the webhook confirms, D3).
func cardIntent() *payssage.Intent {
	return &payssage.Intent{
		ID:       uuid.New(),
		Status:   payssage.IntentStatusSucceeded,
		SellerID: uuid.New(),
	}
}

// ── Input builders ─────────────────────────────────────────────────────────

func ticketLine(ticketID uuid.UUID, attendee *checkouts.Attendee) checkouts.CheckoutLine {
	return checkouts.CheckoutLine{
		ItemType: models.PurchaseItemTypeTicket,
		ItemID:   ticketID,
		Quantity: 1,
		Attendee: attendee,
	}
}

func productLine(variantID uuid.UUID, qty int) checkouts.CheckoutLine {
	return checkouts.CheckoutLine{
		ItemType: models.PurchaseItemTypeProduct,
		ItemID:   variantID,
		Quantity: qty,
	}
}

func programLine(occurrenceID uuid.UUID) checkouts.CheckoutLine {
	return checkouts.CheckoutLine{
		ItemType: models.PurchaseItemTypeProgramOccurrence,
		ItemID:   occurrenceID,
		Quantity: 1,
	}
}

func selfAttendee(userID uuid.UUID) *checkouts.Attendee {
	return &checkouts.Attendee{UserID: userID, Email: "buyer@example.com", Name: "Jane Doe"}
}

func pixInput(lines ...checkouts.CheckoutLine) checkouts.CheckoutInput {
	return checkouts.CheckoutInput{
		PaymentMethod: "pix",
		Payer:         checkouts.Payer{Email: "buyer@example.com", IdentificationType: "CPF", IdentificationNumber: "12345678900"},
		Items:         lines,
	}
}
