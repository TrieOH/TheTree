package ws_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"univents/internal/services/ws"
	"univents/models"

	payssage "sdk/payssage"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func TestMain(m *testing.M) {
	fun.SetConfig(fun.Config{
		DefaultModule:        "test",
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
	})
	os.Exit(m.Run())
}

// ── fakes ────────────────────────────────────────────────────────────────

// fakeTokenStore is an in-memory ws_tokens repo with the one-time semantics
// of the SQL consume (missing/already-used/expired → nil).
type fakeTokenStore struct {
	mu     sync.Mutex
	tokens map[string]*models.WsToken
}

func newFakeTokenStore() *fakeTokenStore {
	return &fakeTokenStore{tokens: make(map[string]*models.WsToken)}
}

func (f *fakeTokenStore) Create(_ context.Context, t *models.WsToken) (*models.WsToken, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	row := *t
	row.ID = uuid.New()
	f.tokens[t.TokenHash] = &row
	return &row, nil
}

func (f *fakeTokenStore) Consume(_ context.Context, hash string) (*models.WsToken, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	t, ok := f.tokens[hash]
	if !ok || t.UsedAt != nil || !t.ExpiresAt.After(time.Now()) {
		//nolint:nilnil // guard miss — the handshake rejects (mirrors the SQL consume)
		return nil, nil
	}
	now := time.Now()
	t.UsedAt = &now
	return t, nil
}

// fakePurchaseRepo is an in-memory ports.PurchaseRepo: purchases + items
// keyed by id, plus a canned availability read for the store tests.
type fakePurchaseRepo struct {
	mu        sync.Mutex
	purchases map[uuid.UUID]*models.Purchase
	items     map[uuid.UUID][]models.PurchaseItem
	avail     []models.ItemAvailability
}

func newFakePurchaseRepo() *fakePurchaseRepo {
	return &fakePurchaseRepo{
		purchases: make(map[uuid.UUID]*models.Purchase),
		items:     make(map[uuid.UUID][]models.PurchaseItem),
	}
}

func (f *fakePurchaseRepo) seed(p *models.Purchase, items []models.PurchaseItem) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p.ID = uuid.New()
	for i := range items {
		items[i].ID = uuid.New()
		items[i].PurchaseID = p.ID
	}
	f.purchases[p.ID] = p
	f.items[p.ID] = items
}

func (f *fakePurchaseRepo) get(id uuid.UUID) *models.Purchase {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.purchases[id]
	if !ok {
		return nil
	}
	cp := *p
	return &cp
}

func (f *fakePurchaseRepo) CreatePurchase(context.Context, *models.Purchase) (*models.Purchase, error) {
	return nil, fun.ErrInternal("unused fake")
}
func (f *fakePurchaseRepo) CreatePurchaseItem(context.Context, *models.PurchaseItem) (*models.PurchaseItem, error) {
	return nil, fun.ErrInternal("unused fake")
}
func (f *fakePurchaseRepo) GetByID(_ context.Context, id uuid.UUID) (*models.Purchase, error) {
	p := f.get(id)
	if p == nil {
		return nil, fun.ErrNotFound("purchase not found")
	}
	return p, nil
}
func (f *fakePurchaseRepo) GetByIDForOwner(_ context.Context, id, purchaserID uuid.UUID) (*models.Purchase, error) {
	p := f.get(id)
	if p == nil || p.PurchaserID != purchaserID {
		return nil, fun.ErrNotFound("purchase not found")
	}
	return p, nil
}
func (f *fakePurchaseRepo) GetByIntentID(context.Context, uuid.UUID) (*models.Purchase, error) {
	return nil, fun.ErrInternal("unused fake")
}
func (f *fakePurchaseRepo) ListByPurchaser(context.Context, uuid.UUID) ([]models.Purchase, error) {
	return nil, fun.ErrInternal("unused fake")
}
func (f *fakePurchaseRepo) UpdateStatus(_ context.Context, id uuid.UUID, status models.PurchaseStatus, reason *string) (*models.Purchase, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.purchases[id]
	if !ok {
		return nil, fun.ErrNotFound("purchase not found")
	}
	p.Status = status
	p.StatusReason = reason
	cp := *p
	return &cp, nil
}
func (f *fakePurchaseRepo) UpdateStatusIf(context.Context, uuid.UUID, models.PurchaseStatus, models.PurchaseStatus, *string) (*models.Purchase, error) {
	return nil, fun.ErrInternal("unused fake")
}
func (f *fakePurchaseRepo) UpdateRiverJobID(context.Context, uuid.UUID, int64) (*models.Purchase, error) {
	return nil, fun.ErrInternal("unused fake")
}
func (f *fakePurchaseRepo) AttachIntent(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, *string, *string) (*models.Purchase, error) {
	return nil, fun.ErrInternal("unused fake")
}
func (f *fakePurchaseRepo) ListItemsByPurchase(_ context.Context, purchaseID uuid.UUID) ([]models.PurchaseItem, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]models.PurchaseItem(nil), f.items[purchaseID]...), nil
}
func (f *fakePurchaseRepo) Availability(context.Context, uuid.UUID) ([]models.ItemAvailability, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]models.ItemAvailability(nil), f.avail...), nil
}

// fakeIntents returns canned GetIntent results.
type fakeIntents struct {
	status payssage.IntentStatus
	detail *string
	err    error
}

func (f *fakeIntents) GetIntent(context.Context, uuid.UUID) (*payssage.Intent, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &payssage.Intent{Status: f.status, StatusDetail: f.detail}, nil
}

// fakeNotifier captures Subscribe handlers; fire() replays a payload to all
// of them (the notifier would deliver via LISTEN/NOTIFY).
type fakeNotifier struct {
	mu       sync.Mutex
	handlers []func(payload string)
}

func newFakeNotifier() *fakeNotifier {
	return &fakeNotifier{}
}

func (f *fakeNotifier) Subscribe(_ string, handler func(payload string)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.handlers = append(f.handlers, handler)
}

func (f *fakeNotifier) fire(payload string) {
	f.mu.Lock()
	handlers := append([]func(string){}, f.handlers...)
	f.mu.Unlock()
	for _, h := range handlers {
		h(payload)
	}
}

// ── fixtures ─────────────────────────────────────────────────────────────

// newOps wires the ws service over the fakes.
func newOps(t *testing.T) (*ws.Operations, *fakeTokenStore, *fakePurchaseRepo, *fakeNotifier) {
	t.Helper()
	return newOpsWithIntents(t, &fakeIntents{status: payssage.IntentStatusPending})
}

// newOpsWithIntents wires the ws service with a caller-provided intent fake.
func newOpsWithIntents(t *testing.T, intents *fakeIntents) (*ws.Operations, *fakeTokenStore, *fakePurchaseRepo, *fakeNotifier) {
	t.Helper()
	tokens := newFakeTokenStore()
	purchases := newFakePurchaseRepo()
	notifier := newFakeNotifier()
	ops := ws.NewOperations(tokens, purchases, intents, notifier)
	ops.SetIntentRetry(1, time.Millisecond)
	return ops, tokens, purchases, notifier
}

// seedPending creates the canonical test purchase: pending, with an intent.
func seedPending(purchases *fakePurchaseRepo) *models.Purchase {
	p := &models.Purchase{
		EditionID:        uuid.New(),
		PurchaserID:      uuid.New(),
		Status:           models.PurchaseStatusPending,
		TotalCents:       1234,
		Currency:         "BRL",
		PayssageIntentID: new(uuid.UUID),
		ExpiresAt:        time.Now().Add(10 * time.Minute),
		CreatedAt:        time.Now(),
	}
	items := []models.PurchaseItem{
		{ItemType: models.PurchaseItemTypeTicket, ItemID: uuid.New(), Quantity: 2, UnitPriceCents: 617},
	}
	purchases.seed(p, items)
	return p
}
