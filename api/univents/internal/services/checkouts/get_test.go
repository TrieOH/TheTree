package checkouts_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	"lib/database"
	"lib/testdb"

	"univents/internal/repos"
	"univents/internal/services/checkouts"
	"univents/internal/sqlc"
	"univents/models"

	"sdk/payssage"
)

// intentStub fakes the IntentClient seam: each GetIntent call runs fn with
// the 1-based attempt number, so tests can simulate transient failures
// (fail twice, succeed on the third) and persistent failure.
type intentStub struct {
	mu    sync.Mutex
	calls int
	fn    func(intentID uuid.UUID, attempt int) (*payssage.Intent, error)
}

func (s *intentStub) GetIntent(_ context.Context, intentID uuid.UUID) (*payssage.Intent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	return s.fn(intentID, s.calls)
}

func (s *intentStub) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

// seedEdition creates an event + edition the purchases can reference.
func seedEdition(t *testing.T, q *sqlc.Queries) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	event, err := q.CreateEvent(ctx, sqlc.CreateEventParams{
		OwnerID:  uuid.New(),
		FullName: "Read Surfaces Test Event",
		Slug:     "read-surfaces-" + uuid.NewString()[:8],
		Status:   string(models.EventStatusActive),
	})
	if err != nil {
		t.Fatalf("seed event: %v", err)
	}
	edition, err := q.CreateEdition(ctx, sqlc.CreateEditionParams{
		EventID:     event.ID,
		EditionName: "Read Surfaces Test Edition",
		Slug:        "read-surfaces-ed-" + uuid.NewString()[:8],
		StartsAt:    time.Now().Add(-time.Hour),
		EndsAt:      time.Now().Add(24 * time.Hour),
		CreatedBy:   event.OwnerID,
	})
	if err != nil {
		t.Fatalf("seed edition: %v", err)
	}
	return edition.ID
}

// seedPurchase creates one purchase with one item of each type (the D4
// ledger shape a split-7 checkout would produce) and returns it with its
// items.
func seedPurchase(t *testing.T, r *repos.Repos, editionID, purchaserID uuid.UUID, status models.PurchaseStatus, intentID *uuid.UUID) (*models.Purchase, []models.PurchaseItem) {
	t.Helper()
	ctx := context.Background()

	purchase, err := r.Purchases.CreatePurchase(ctx, &models.Purchase{
		EditionID:        editionID,
		PurchaserID:      purchaserID,
		Status:           status,
		TotalCents:       8000,
		Currency:         "BRL",
		PaymentMethod:    new("pix"),
		PayssageIntentID: intentID,
		ExpiresAt:        time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		t.Fatalf("seed purchase: %v", err)
	}

	items := []models.PurchaseItem{
		{ItemType: models.PurchaseItemTypeTicket, ItemID: uuid.New(), Quantity: 1, UnitPriceCents: 1000},
		{ItemType: models.PurchaseItemTypeProduct, ItemID: uuid.New(), Quantity: 2, UnitPriceCents: 5000},
		{ItemType: models.PurchaseItemTypeProgramOccurrence, ItemID: uuid.New(), Quantity: 1, UnitPriceCents: 2000},
	}
	for i := range items {
		items[i].PurchaseID = purchase.ID
		created, err := r.Purchases.CreatePurchaseItem(ctx, &items[i])
		if err != nil {
			t.Fatalf("seed purchase item: %v", err)
		}
		_ = created
	}
	return purchase, items
}

func alwaysSucceed(_ uuid.UUID, _ int) (*payssage.Intent, error) {
	return &payssage.Intent{Status: payssage.IntentStatusSucceeded}, nil
}

// newOps wires the real repos (disposable Postgres with the real
// migrations) behind a faked IntentClient. The intent retry delay is zeroed
// so retry tests stay instant.
func newOps(t *testing.T, intents checkouts.IntentClient) (*repos.Repos, *sqlc.Queries, *checkouts.Operations) {
	t.Helper()
	pool := testdb.Postgres(t, "../../../db/migrations")
	q := sqlc.New(pool)
	database.SetDefaultRunner(database.NewPGXTxRunner(pool))
	r := repos.New(q)

	ops := checkouts.NewOperations(r.Purchases, intents)
	ops.SetIntentRetry(3, 0)
	return r, q, ops
}

func TestGet_OwnerSeesFullState(t *testing.T) {
	intents := &intentStub{fn: alwaysSucceed}
	r, q, ops := newOps(t, intents)
	editionID := seedEdition(t, q)
	purchaserID := uuid.New()
	intentID := uuid.New()
	purchase, seededItems := seedPurchase(t, r, editionID, purchaserID, models.PurchaseStatusPending, &intentID)

	resume, err := ops.Get(context.Background(), purchase.ID, purchaserID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if resume.Purchase.ID != purchase.ID {
		t.Fatalf("purchase id = %s, want %s", resume.Purchase.ID, purchase.ID)
	}
	if resume.Purchase.Status != models.PurchaseStatusPending {
		t.Fatalf("status = %s, want pending", resume.Purchase.Status)
	}
	if resume.Purchase.TotalCents != 8000 || resume.Purchase.Currency != "BRL" {
		t.Fatalf("total/currency = %d %s, want 8000 BRL", resume.Purchase.TotalCents, resume.Purchase.Currency)
	}
	if resume.Purchase.PayssageIntentID == nil || *resume.Purchase.PayssageIntentID != intentID {
		t.Fatalf("intent id = %v, want %s", resume.Purchase.PayssageIntentID, intentID)
	}

	if len(resume.Items) != len(seededItems) {
		t.Fatalf("items = %d, want %d", len(resume.Items), len(seededItems))
	}
	for i, item := range resume.Items {
		want := seededItems[i]
		if item.ItemType != want.ItemType || item.ItemID != want.ItemID ||
			item.Quantity != want.Quantity || item.UnitPriceCents != want.UnitPriceCents {
			t.Fatalf("item %d = %+v, want %+v", i, item, want)
		}
	}

	if resume.IntentStatus == nil || *resume.IntentStatus != string(payssage.IntentStatusSucceeded) {
		t.Fatalf("intent_status = %v, want succeeded", resume.IntentStatus)
	}
}

func TestGet_NonOwnerIsNotFound(t *testing.T) {
	r, q, ops := newOps(t, &intentStub{fn: alwaysSucceed})
	editionID := seedEdition(t, q)
	purchase, _ := seedPurchase(t, r, editionID, uuid.New(), models.PurchaseStatusPending, new(uuid.UUID))

	// A different authenticated user must not see the purchase — 404 (no
	// existence leak), never 403.
	_, err := ops.Get(context.Background(), purchase.ID, uuid.New())
	if err == nil {
		t.Fatal("Get: expected error for non-owner, got nil")
	}
	if !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("Get: error = %v, want NOT_FOUND", err)
	}
}

func TestGet_UnknownPurchaseIsNotFound(t *testing.T) {
	_, _, ops := newOps(t, &intentStub{fn: alwaysSucceed})
	_, err := ops.Get(context.Background(), uuid.New(), uuid.New())
	if err == nil || !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("Get: err = %v, want NOT_FOUND", err)
	}
}

// TestGet_IntentFetchRetry pins the retry budget: transient failures (two
// misses) still surface the intent status on the third attempt, and the
// intent seam is called exactly 3 times.
func TestGet_IntentFetchRetry(t *testing.T) {
	intents := &intentStub{fn: func(_ uuid.UUID, attempt int) (*payssage.Intent, error) {
		if attempt < 3 {
			return nil, errors.New("payssage temporarily unavailable")
		}
		return &payssage.Intent{Status: payssage.IntentStatusSucceeded}, nil
	}}
	r, q, ops := newOps(t, intents)
	editionID := seedEdition(t, q)
	purchaserID := uuid.New()
	intentID := uuid.New()
	purchase, _ := seedPurchase(t, r, editionID, purchaserID, models.PurchaseStatusPending, &intentID)

	resume, err := ops.Get(context.Background(), purchase.ID, purchaserID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if resume.IntentStatus == nil || *resume.IntentStatus != string(payssage.IntentStatusSucceeded) {
		t.Fatalf("intent_status = %v, want succeeded after retries", resume.IntentStatus)
	}
	if got := intents.callCount(); got != 3 {
		t.Fatalf("GetIntent calls = %d, want 3", got)
	}
}

// TestGet_IntentFetchDegrades pins the graceful-degrade path: when the
// intent never resolves, the resume still returns the purchase's own state
// (source of truth) with a nil intent_status — no error, no 500.
func TestGet_IntentFetchDegrades(t *testing.T) {
	intents := &intentStub{fn: func(_ uuid.UUID, _ int) (*payssage.Intent, error) {
		return nil, errors.New("payssage unreachable")
	}}
	r, q, ops := newOps(t, intents)
	editionID := seedEdition(t, q)
	purchaserID := uuid.New()
	intentID := uuid.New()
	purchase, _ := seedPurchase(t, r, editionID, purchaserID, models.PurchaseStatusPending, &intentID)

	resume, err := ops.Get(context.Background(), purchase.ID, purchaserID)
	if err != nil {
		t.Fatalf("Get: %v (must degrade, not fail)", err)
	}
	if resume.IntentStatus != nil {
		t.Fatalf("intent_status = %v, want nil (degraded)", resume.IntentStatus)
	}
	if got := intents.callCount(); got != 3 {
		t.Fatalf("GetIntent calls = %d, want 3 (full budget before degrading)", got)
	}
}

// TestGet_NoIntentFetchWhenNotPending pins the fetch guard: the intent is
// only consulted when the purchase is pending — an approved purchase never
// touches Payssage on resume.
func TestGet_NoIntentFetchWhenNotPending(t *testing.T) {
	intents := &intentStub{fn: alwaysSucceed}
	r, q, ops := newOps(t, intents)
	editionID := seedEdition(t, q)
	purchaserID := uuid.New()
	intentID := uuid.New()
	purchase, _ := seedPurchase(t, r, editionID, purchaserID, models.PurchaseStatusApproved, &intentID)

	resume, err := ops.Get(context.Background(), purchase.ID, purchaserID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if resume.IntentStatus != nil {
		t.Fatalf("intent_status = %v, want nil (not pending)", resume.IntentStatus)
	}
	if got := intents.callCount(); got != 0 {
		t.Fatalf("GetIntent calls = %d, want 0", got)
	}
}

// TestGet_NoIntentFetchWithoutIntent pins the free-order case: a pending
// purchase without a Payssage intent (total 0, split 7) has nothing to
// fetch.
func TestGet_NoIntentFetchWithoutIntent(t *testing.T) {
	intents := &intentStub{fn: alwaysSucceed}
	r, q, ops := newOps(t, intents)
	editionID := seedEdition(t, q)
	purchaserID := uuid.New()
	purchase, _ := seedPurchase(t, r, editionID, purchaserID, models.PurchaseStatusPending, nil)

	resume, err := ops.Get(context.Background(), purchase.ID, purchaserID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if resume.IntentStatus != nil {
		t.Fatalf("intent_status = %v, want nil (no intent)", resume.IntentStatus)
	}
	if got := intents.callCount(); got != 0 {
		t.Fatalf("GetIntent calls = %d, want 0", got)
	}
}
