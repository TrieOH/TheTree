package webhooks_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"lib/database"

	"univents/internal/services/webhooks"
	"univents/models"
	"univents/ports"
)

func TestMain(m *testing.M) {
	fun.SetConfig(fun.Config{
		DefaultModule:        "test",
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
	})
	os.Exit(m.Run())
}

const testSecret = "webhook-test-secret" //nolint:gosec // test-only shared secret, not a credential

// sign reproduces payssage's X-Payssage-Signature over the raw body.
func sign(body []byte) string {
	mac := hmac.New(sha256.New, []byte(testSecret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func input(intentID uuid.UUID, eventType string, body []byte) webhooks.ReceiveInput {
	return webhooks.ReceiveInput{
		IntentID:  intentID,
		EventType: eventType,
		RawBody:   body,
		Signature: sign(body),
	}
}

// noopTxRunner runs fn directly — the repo fakes ignore the tx context, so
// the tx is simulated. The DB-backed test exercises the real pgx runner.
type noopTxRunner struct{}

func (noopTxRunner) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
func (noopTxRunner) WithinTxWithOptions(ctx context.Context, _ database.TxOptions, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type notifyCall struct {
	channel string
	payload string
}

// recordingNotifier captures Notify calls for assertions.
type recordingNotifier struct {
	mu    sync.Mutex
	calls []notifyCall
}

func (n *recordingNotifier) Notify(_ context.Context, channel, payload string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.calls = append(n.calls, notifyCall{channel, payload})
	return nil
}

func (n *recordingNotifier) payloads() []string {
	n.mu.Lock()
	defer n.mu.Unlock()
	out := make([]string, 0, len(n.calls))
	for _, c := range n.calls {
		out = append(out, c.payload)
	}
	return out
}

// decodedNotify mirrors the notifier payload schema (split 6): kind stock or
// purchase.
type decodedNotify struct {
	Kind       string      `json:"kind"`
	EditionID  uuid.UUID   `json:"edition_id"`
	ItemIDs    []uuid.UUID `json:"item_ids"`
	PurchaseID uuid.UUID   `json:"purchase_id"`
	Status     string      `json:"status"`
}

func (n *recordingNotifier) decoded() []decodedNotify {
	out := make([]decodedNotify, 0, len(n.payloads()))
	for _, p := range n.payloads() {
		var d decodedNotify
		err := json.Unmarshal([]byte(p), &d)
		if err != nil {
			panic("bad notify payload: " + err.Error())
		}
		out = append(out, d)
	}
	return out
}

func purchase(id, intentID uuid.UUID, status models.PurchaseStatus) *models.Purchase {
	return &models.Purchase{
		ID:               id,
		EditionID:        uuid.New(),
		PurchaserID:      uuid.New(),
		Status:           status,
		TotalCents:       8000,
		Currency:         "BRL",
		PayssageIntentID: &intentID,
		ExpiresAt:        time.Now().Add(10 * time.Minute),
	}
}

type harness struct {
	purchases        ports.PurchaseRepo
	registrations    ports.RegistrationRepo
	productPurchases ports.ProductPurchaseRepo
	participations   ports.ProgramParticipationRepo
	badges           webhooks.Badges
	river            webhooks.River
	notifier         *recordingNotifier
	ops              *webhooks.Operations
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	mock.SetUp(t)
	h := &harness{
		purchases:        mock.Mock[ports.PurchaseRepo](),
		registrations:    mock.Mock[ports.RegistrationRepo](),
		productPurchases: mock.Mock[ports.ProductPurchaseRepo](),
		participations:   mock.Mock[ports.ProgramParticipationRepo](),
		badges:           mock.Mock[webhooks.Badges](),
		river:            mock.Mock[webhooks.River](),
		notifier:         &recordingNotifier{},
	}
	h.ops = webhooks.NewOperations(
		h.purchases, h.registrations, h.productPurchases, h.participations,
		h.badges, h.notifier, h.river, noopTxRunner{}, testSecret,
	)
	h.ops.SetCardRaceWait(0)
	return h
}

// stubApprove stubs the happy-path approval flow for a purchase with one
// ticket item (registration + badge) and one product item.
func stubApprove(h *harness, p *models.Purchase, items []models.PurchaseItem) {
	approved := *p
	approved.Status = models.PurchaseStatusApproved

	mock.When(h.purchases.ListItemsByPurchase(mock.AnyContext(), mock.Equal(p.ID))).ThenReturn(items, nil)
	mock.When(h.purchases.UpdateStatusIf(mock.AnyContext(), mock.Equal(p.ID), mock.Equal(p.Status), mock.Equal(models.PurchaseStatusApproved), mock.Any[*string]())).
		ThenReturn(&approved, nil)

	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				continue
			}
			mock.When(h.registrations.UpdateStatus(mock.AnyContext(), mock.Equal(*item.RegistrationID), mock.Equal(models.RegistrationStatusConfirmed), mock.Any[*string]())).
				ThenReturn(&models.Registration{ID: *item.RegistrationID, Status: models.RegistrationStatusConfirmed}, nil)
			mock.When(h.badges.EmitForConfirmedRegistration(mock.AnyContext(), mock.Equal(*item.RegistrationID))).
				ThenReturn(&models.BadgeEmission{ID: uuid.New()}, nil)
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				continue
			}
			mock.When(h.productPurchases.UpdateProductPurchaseStatus(mock.AnyContext(), mock.Equal(*item.ProductPurchaseID), mock.Equal(models.ProductPurchaseStatusConfirmed), mock.Any[*string]())).
				ThenReturn(&models.ProductPurchase{ID: *item.ProductPurchaseID, Status: models.ProductPurchaseStatusConfirmed}, nil)
		}
	}
}

func ticketItem(purchaseID, regID uuid.UUID) models.PurchaseItem {
	return models.PurchaseItem{ID: uuid.New(), PurchaseID: purchaseID, ItemType: models.PurchaseItemTypeTicket, ItemID: uuid.New(), Quantity: 1, RegistrationID: &regID}
}

func productItem(purchaseID, ppID uuid.UUID) models.PurchaseItem {
	return models.PurchaseItem{ID: uuid.New(), PurchaseID: purchaseID, ItemType: models.PurchaseItemTypeProduct, ItemID: uuid.New(), Quantity: 1, ProductPurchaseID: &ppID}
}

func TestReceive_RejectsBadSignature(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()

	err := h.ops.Receive(context.Background(), webhooks.ReceiveInput{
		IntentID:  intentID,
		EventType: "payment.succeeded",
		RawBody:   []byte(`{"intent_id":"` + intentID.String() + `"}`),
		Signature: "deadbeef",
	})
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want bad request, got %v", err)
	}
	// No correlation attempt happens before verification.
	_, _ = mock.Verify(h.purchases, mock.Never()).GetByIntentID(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestReceive_ApproveHappyPath(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()
	purchaseID := uuid.New()
	regID := uuid.New()
	ppID := uuid.New()
	p := purchase(purchaseID, intentID, models.PurchaseStatusPending)
	jobID := int64(42)
	p.RiverJobID = &jobID
	items := []models.PurchaseItem{ticketItem(purchaseID, regID), productItem(purchaseID, ppID)}

	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).ThenReturn(p, nil)
	mock.When(h.river.JobCancel(mock.AnyContext(), mock.Equal(jobID))).ThenReturn(nil, nil)
	stubApprove(h, p, items)

	err := h.ops.Receive(context.Background(), input(intentID, "payment.succeeded", []byte(`{}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}

	// Materialized flips + badge emit happened.
	_, _ = mock.Verify(h.registrations, mock.Once()).UpdateStatus(mock.AnyContext(), mock.Equal(regID), mock.Equal(models.RegistrationStatusConfirmed), mock.Any[*string]())
	_, _ = mock.Verify(h.badges, mock.Once()).EmitForConfirmedRegistration(mock.AnyContext(), mock.Equal(regID))
	_, _ = mock.Verify(h.productPurchases, mock.Once()).UpdateProductPurchaseStatus(mock.AnyContext(), mock.Equal(ppID), mock.Equal(models.ProductPurchaseStatusConfirmed), mock.Any[*string]())
	// Expiry job cancelled (best-effort).
	_, _ = mock.Verify(h.river, mock.Once()).JobCancel(mock.AnyContext(), mock.Equal(jobID))
	// Notified: stock deltas + purchase.confirmed.
	notifies := h.notifier.decoded()
	if len(notifies) != 2 {
		t.Fatalf("notifications = %d, want 2 (stock + purchase)", len(notifies))
	}
	byKind := map[string]decodedNotify{}
	for _, n := range notifies {
		byKind[n.Kind] = n
	}
	if stock, ok := byKind["stock"]; !ok || stock.EditionID != p.EditionID || len(stock.ItemIDs) != 2 {
		t.Fatalf("stock notify = %+v, want edition %s + 2 item ids", byKind["stock"], p.EditionID)
	}
	if purch, ok := byKind["purchase"]; !ok || purch.PurchaseID != purchaseID || purch.Status != "approved" {
		t.Fatalf("purchase notify = %+v, want purchase.confirmed", byKind["purchase"])
	}
}

func TestReceive_DuplicateApprovalIsNoOp(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()
	p := purchase(uuid.New(), intentID, models.PurchaseStatusApproved)

	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).ThenReturn(p, nil)

	err := h.ops.Receive(context.Background(), input(intentID, "payment.succeeded", []byte(`{}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}
	_, _ = mock.Verify(h.purchases, mock.Never()).UpdateStatusIf(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[models.PurchaseStatus](), mock.Any[models.PurchaseStatus](), mock.Any[*string]())
	_, _ = mock.Verify(h.registrations, mock.Never()).UpdateStatus(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[models.RegistrationStatus](), mock.Any[*string]())
	if got := h.notifier.payloads(); len(got) != 0 {
		t.Fatalf("duplicate delivery published %d notifications", len(got))
	}
}

func TestReceive_CardRaceResolvesOnRequery(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()
	purchaseID := uuid.New()
	regID := uuid.New()
	p := purchase(purchaseID, intentID, models.PurchaseStatusPending)
	items := []models.PurchaseItem{ticketItem(purchaseID, regID)}

	// First lookup misses (checkout tx not committed yet, D3) — the second
	// (after the ~1s wait, zeroed in tests) finds it.
	calls := 0
	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).
		ThenAnswer(func(_ []any) []any {
			calls++
			if calls == 1 {
				return []any{nil, fun.ErrNotFound("purchase not found")}
			}
			return []any{p, nil}
		})
	stubApprove(h, p, items)

	err := h.ops.Receive(context.Background(), input(intentID, "payment.succeeded", []byte(`{}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}
	_, _ = mock.Verify(h.purchases, mock.Times(2)).GetByIntentID(mock.AnyContext(), mock.Equal(intentID))
	_, _ = mock.Verify(h.badges, mock.Once()).EmitForConfirmedRegistration(mock.AnyContext(), mock.Equal(regID))
}

func TestReceive_CardRaceStillMissingReturnsNon2xx(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()

	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).
		ThenReturn(nil, fun.ErrNotFound("purchase not found"))

	err := h.ops.Receive(context.Background(), input(intentID, "payment.succeeded", []byte(`{}`)))
	if err == nil || !fun.Is(err, fun.CodeInternal) {
		t.Fatalf("want internal (non-2xx → payssage retries), got %v", err)
	}
	_, _ = mock.Verify(h.purchases, mock.Times(2)).GetByIntentID(mock.AnyContext(), mock.Equal(intentID))
}

func TestReceive_LateApprovalHonoredWithFullStock(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()
	purchaseID := uuid.New()
	regID := uuid.New()
	itemID := uuid.New()
	p := purchase(purchaseID, intentID, models.PurchaseStatusExpired)
	items := []models.PurchaseItem{{
		ID: uuid.New(), PurchaseID: purchaseID, ItemType: models.PurchaseItemTypeTicket,
		ItemID: itemID, Quantity: 2, RegistrationID: &regID,
	}}

	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).ThenReturn(p, nil)
	mock.When(h.purchases.Availability(mock.AnyContext(), mock.Equal(p.EditionID))).
		ThenReturn([]models.ItemAvailability{{ItemType: models.PurchaseItemTypeTicket, ItemID: itemID, BaseQuantity: new(int(10)), ReservedQuantity: 3}}, nil)
	stubApprove(h, p, items)

	err := h.ops.Receive(context.Background(), input(intentID, "payment.approved", []byte(`{}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}
	_, _ = mock.Verify(h.purchases, mock.Once()).UpdateStatusIf(mock.AnyContext(), mock.Equal(purchaseID), mock.Equal(models.PurchaseStatusExpired), mock.Equal(models.PurchaseStatusApproved), mock.Any[*string]())
	_, _ = mock.Verify(h.badges, mock.Once()).EmitForConfirmedRegistration(mock.AnyContext(), mock.Equal(regID))
}

func TestReceive_LateApprovalWithoutStockDefersRefund(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()
	purchaseID := uuid.New()
	regID := uuid.New()
	itemID := uuid.New()
	p := purchase(purchaseID, intentID, models.PurchaseStatusExpired)
	items := []models.PurchaseItem{{
		ID: uuid.New(), PurchaseID: purchaseID, ItemType: models.PurchaseItemTypeTicket,
		ItemID: itemID, Quantity: 2, RegistrationID: &regID,
	}}

	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).ThenReturn(p, nil)
	mock.When(h.purchases.ListItemsByPurchase(mock.AnyContext(), mock.Equal(purchaseID))).ThenReturn(items, nil)
	mock.When(h.purchases.Availability(mock.AnyContext(), mock.Equal(p.EditionID))).
		ThenReturn([]models.ItemAvailability{{ItemType: models.PurchaseItemTypeTicket, ItemID: itemID, BaseQuantity: new(int(2)), ReservedQuantity: 2}}, nil)

	var gotReason *string
	mock.When(h.purchases.UpdateStatus(mock.AnyContext(), mock.Equal(purchaseID), mock.Equal(models.PurchaseStatusExpired), mock.Any[*string]())).
		ThenAnswer(func(args []any) []any {
			gotReason = args[3].(*string)
			return []any{p, nil}
		})

	err := h.ops.Receive(context.Background(), input(intentID, "payment.succeeded", []byte(`{}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}
	// Refund deferred (D1): stays expired, reason set, no approval, no badge.
	if gotReason == nil || *gotReason != "refunded_after_expiry" {
		t.Fatalf("status_reason = %v, want refunded_after_expiry", gotReason)
	}
	_, _ = mock.Verify(h.purchases, mock.Never()).UpdateStatusIf(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[models.PurchaseStatus](), mock.Any[models.PurchaseStatus](), mock.Any[*string]())
	_, _ = mock.Verify(h.badges, mock.Never()).EmitForConfirmedRegistration(mock.AnyContext(), mock.Any[uuid.UUID]())
	// Notified with the expired purchase event.
	notifies := h.notifier.decoded()
	statuses := map[string]bool{}
	for _, n := range notifies {
		statuses[n.Kind] = true
		if n.Kind == "purchase" && n.Status != "expired" {
			t.Fatalf("purchase notify status = %s, want expired", n.Status)
		}
	}
	if !statuses["stock"] || !statuses["purchase"] {
		t.Fatalf("want stock + purchase notifications, got %v", notifies)
	}
}

func TestReceive_FailureCancelsPendingPurchase(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()
	purchaseID := uuid.New()
	regID := uuid.New()
	ppID := uuid.New()
	partID := uuid.New()
	p := purchase(purchaseID, intentID, models.PurchaseStatusPending)
	items := []models.PurchaseItem{
		ticketItem(purchaseID, regID),
		productItem(purchaseID, ppID),
		{ID: uuid.New(), PurchaseID: purchaseID, ItemType: models.PurchaseItemTypeProgramOccurrence, ItemID: uuid.New(), Quantity: 1, ParticipationID: &partID},
	}
	cancelled := *p
	cancelled.Status = models.PurchaseStatusCancelled

	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).ThenReturn(p, nil)
	mock.When(h.purchases.ListItemsByPurchase(mock.AnyContext(), mock.Equal(purchaseID))).ThenReturn(items, nil)
	mock.When(h.purchases.UpdateStatusIf(mock.AnyContext(), mock.Equal(purchaseID), mock.Equal(models.PurchaseStatusPending), mock.Equal(models.PurchaseStatusCancelled), mock.Any[*string]())).
		ThenReturn(&cancelled, nil)
	mock.When(h.registrations.UpdateStatus(mock.AnyContext(), mock.Equal(regID), mock.Equal(models.RegistrationStatusCancelled), mock.Any[*string]())).
		ThenReturn(&models.Registration{ID: regID, Status: models.RegistrationStatusCancelled}, nil)
	mock.When(h.productPurchases.UpdateProductPurchaseStatus(mock.AnyContext(), mock.Equal(ppID), mock.Equal(models.ProductPurchaseStatusCancelled), mock.Any[*string]())).
		ThenReturn(&models.ProductPurchase{ID: ppID, Status: models.ProductPurchaseStatusCancelled}, nil)
	mock.When(h.participations.UpdateParticipationStatus(mock.AnyContext(), mock.Equal(partID), mock.Equal(models.ProgramParticipationStatusCancelled))).
		ThenReturn(&models.ProgramParticipation{ID: partID, Status: models.ProgramParticipationStatusCancelled}, nil)

	err := h.ops.Receive(context.Background(), input(intentID, "payment.rejected", []byte(`{}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}
	_, _ = mock.Verify(h.registrations, mock.Once()).UpdateStatus(mock.AnyContext(), mock.Equal(regID), mock.Equal(models.RegistrationStatusCancelled), mock.Any[*string]())
	_, _ = mock.Verify(h.productPurchases, mock.Once()).UpdateProductPurchaseStatus(mock.AnyContext(), mock.Equal(ppID), mock.Equal(models.ProductPurchaseStatusCancelled), mock.Any[*string]())
	_, _ = mock.Verify(h.participations, mock.Once()).UpdateParticipationStatus(mock.AnyContext(), mock.Equal(partID), mock.Equal(models.ProgramParticipationStatusCancelled))
	// No badge emit on failure.
	_, _ = mock.Verify(h.badges, mock.Never()).EmitForConfirmedRegistration(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestReceive_FailureForNonPendingPurchaseIgnored(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()
	p := purchase(uuid.New(), intentID, models.PurchaseStatusApproved)

	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).ThenReturn(p, nil)

	err := h.ops.Receive(context.Background(), input(intentID, "payment.cancelled", []byte(`{}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}
	_, _ = mock.Verify(h.purchases, mock.Never()).UpdateStatusIf(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[models.PurchaseStatus](), mock.Any[models.PurchaseStatus](), mock.Any[*string]())
	if got := h.notifier.payloads(); len(got) != 0 {
		t.Fatalf("stale failure published %d notifications", len(got))
	}
}

func TestReceive_RefundedIgnored(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()
	p := purchase(uuid.New(), intentID, models.PurchaseStatusApproved)

	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).ThenReturn(p, nil)

	err := h.ops.Receive(context.Background(), input(intentID, "payment.refunded", []byte(`{}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}
	_, _ = mock.Verify(h.purchases, mock.Never()).UpdateStatusIf(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[models.PurchaseStatus](), mock.Any[models.PurchaseStatus](), mock.Any[*string]())
	if got := h.notifier.payloads(); len(got) != 0 {
		t.Fatalf("refunded published %d notifications", len(got))
	}
}

func TestReceive_PendingEventIsAck(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()
	p := purchase(uuid.New(), intentID, models.PurchaseStatusPending)

	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).ThenReturn(p, nil)

	err := h.ops.Receive(context.Background(), input(intentID, "payment.pending", []byte(`{}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}
	_, _ = mock.Verify(h.purchases, mock.Never()).UpdateStatusIf(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[models.PurchaseStatus](), mock.Any[models.PurchaseStatus](), mock.Any[*string]())
	if got := h.notifier.payloads(); len(got) != 0 {
		t.Fatalf("pending ack published %d notifications", len(got))
	}
}

func TestReceive_UnhandledEventTypeAcked(t *testing.T) {
	h := newHarness(t)
	intentID := uuid.New()
	p := purchase(uuid.New(), intentID, models.PurchaseStatusPending)

	mock.When(h.purchases.GetByIntentID(mock.AnyContext(), mock.Equal(intentID))).ThenReturn(p, nil)

	err := h.ops.Receive(context.Background(), input(intentID, "payment.whatsapp", []byte(`{}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}
	if got := h.notifier.payloads(); len(got) != 0 {
		t.Fatalf("unknown event published %d notifications", len(got))
	}
}
