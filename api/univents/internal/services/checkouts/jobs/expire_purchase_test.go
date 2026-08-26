package jobs_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"

	"lib/database"
	"lib/testdb"

	"univents/internal/repos"
	"univents/internal/services/checkouts/jobs"
	"univents/internal/sqlc"
	"univents/models"
)

// recordingNotifier records NOTIFY payloads (stock deltas + purchase event).
type recordingNotifier struct {
	mu       sync.Mutex
	payloads []string
}

func (n *recordingNotifier) Notify(_ context.Context, _ string, payload string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.payloads = append(n.payloads, payload)
	return nil
}

func (n *recordingNotifier) count() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.payloads)
}

// seedExpiring seeds a pending purchase with one materialized row of each
// type (the split-7 checkout shape) and returns it with its items.
func seedExpiring(t *testing.T, r *repos.Repos, status models.PurchaseStatus) (*models.Purchase, []models.PurchaseItem) {
	t.Helper()
	ctx := context.Background()

	event, err := r.Events.Create(ctx, &models.Event{
		OwnerID:  uuid.New(),
		FullName: "Expiry Test Event",
		Slug:     "expiry-test-" + uuid.NewString()[:8],
		Status:   models.EventStatusActive,
	})
	if err != nil {
		t.Fatalf("seed event: %v", err)
	}
	edition, err := r.Editions.Create(ctx, &models.Edition{
		EventID:   event.ID,
		Name:      "Expiry Test Edition",
		Slug:      "expiry-test-ed-" + uuid.NewString()[:8],
		StartsAt:  time.Now().Add(-time.Hour),
		EndsAt:    time.Now().Add(24 * time.Hour),
		CreatedBy: event.OwnerID,
	})
	if err != nil {
		t.Fatalf("seed edition: %v", err)
	}

	ticket, err := r.TicketTypes.Create(ctx, &models.TicketType{
		EditionID:   edition.ID,
		Name:        "Standard",
		AccessLevel: 0,
		PriceCents:  1000,
		MaxQuantity: new(int(10)),
	})
	if err != nil {
		t.Fatalf("seed ticket: %v", err)
	}
	product, err := r.Products.CreateProduct(ctx, &models.Product{
		EditionID:  edition.ID,
		VendorCode: "P-" + uuid.NewString()[:8],
	})
	if err != nil {
		t.Fatalf("seed product: %v", err)
	}
	variant, err := r.Products.CreateVariant(ctx, &models.ProductVariant{
		EditionID:  edition.ID,
		ProductID:  product.ID,
		VendorCode: "V-" + uuid.NewString()[:8],
		Name:       "T-Shirt",
		Price:      5000,
		Stock:      new(int(5)),
	})
	if err != nil {
		t.Fatalf("seed variant: %v", err)
	}
	program, err := r.Programs.Create(ctx, &models.Program{
		EditionID: edition.ID,
		Kind:      models.ProgramKindActivity,
		Name:      "Workshop",
		Price:     new(int64(2000)),
	})
	if err != nil {
		t.Fatalf("seed program: %v", err)
	}
	occurrence, err := r.Programs.CreateOccurrence(ctx, &models.ProgramOccurrence{
		ProgramID:   program.ID,
		EditionID:   edition.ID,
		StartsAt:    time.Now().Add(2 * time.Hour),
		EndsAt:      time.Now().Add(3 * time.Hour),
		MaxCapacity: new(int(3)),
	})
	if err != nil {
		t.Fatalf("seed occurrence: %v", err)
	}

	purchaserID := uuid.New()
	jobID := int64(1)
	purchase, err := r.Purchases.CreatePurchase(ctx, &models.Purchase{
		EditionID:   edition.ID,
		PurchaserID: purchaserID,
		Status:      status,
		TotalCents:  8000,
		Currency:    "BRL",
		ExpiresAt:   time.Now().Add(10 * time.Minute),
		RiverJobID:  &jobID,
	})
	if err != nil {
		t.Fatalf("seed purchase: %v", err)
	}

	reg, err := r.Registrations.Create(ctx, &models.Registration{
		EditionID:      edition.ID,
		TicketTypeID:   ticket.ID,
		PurchaserID:    purchaserID,
		AttendeeUserID: &purchaserID,
		AttendeeEmail:  "buyer@example.com",
		AttendeeName:   "Jane Doe",
		Status:         models.RegistrationStatusPending,
	})
	if err != nil {
		t.Fatalf("seed registration: %v", err)
	}
	pp, err := r.Products.CreateProductPurchase(ctx, &models.ProductPurchase{
		EditionID:      edition.ID,
		VariantID:      variant.ID,
		PurchaserID:    purchaserID,
		RegistrationID: &reg.ID,
		Quantity:       1,
		Status:         models.ProductPurchaseStatusPending,
	})
	if err != nil {
		t.Fatalf("seed product purchase: %v", err)
	}
	part, err := r.Programs.CreateParticipation(ctx, &models.ProgramParticipation{
		EditionID:      edition.ID,
		OccurrenceID:   occurrence.ID,
		RegistrationID: reg.ID,
		Status:         models.ProgramParticipationStatusRegistered,
	})
	if err != nil {
		t.Fatalf("seed participation: %v", err)
	}

	items := []models.PurchaseItem{
		{ItemType: models.PurchaseItemTypeTicket, ItemID: ticket.ID, Quantity: 1, UnitPriceCents: 1000, RegistrationID: &reg.ID},
		{ItemType: models.PurchaseItemTypeProduct, ItemID: variant.ID, Quantity: 1, UnitPriceCents: 5000, ProductPurchaseID: &pp.ID},
		{ItemType: models.PurchaseItemTypeProgramOccurrence, ItemID: occurrence.ID, Quantity: 1, UnitPriceCents: 2000, ParticipationID: &part.ID},
	}
	for i := range items {
		items[i].PurchaseID = purchase.ID
		_, err := r.Purchases.CreatePurchaseItem(ctx, &items[i])
		if err != nil {
			t.Fatalf("seed purchase item: %v", err)
		}
	}
	return purchase, items
}

func newWorker(t *testing.T) (*repos.Repos, *jobs.ExpirePurchaseWorker, *recordingNotifier) {
	t.Helper()
	pool := testdb.Postgres(t, "../../../../db/migrations")
	q := sqlc.New(pool)
	tx := database.NewPGXTxRunner(pool)
	database.SetDefaultRunner(tx)
	r := repos.New(q)
	notifier := &recordingNotifier{}
	w := jobs.NewExpirePurchaseWorker(r.Purchases, r.Registrations, r.Products, r.Programs, notifier, tx)
	return r, w, notifier
}

func run(t *testing.T, w *jobs.ExpirePurchaseWorker, purchaseID uuid.UUID) error {
	t.Helper()
	job := &river.Job[jobs.ExpirePurchaseArgs]{
		Args: jobs.ExpirePurchaseArgs{PurchaseID: purchaseID},
	}
	return w.Work(context.Background(), job)
}

// TestExpire_PendingExpires pins the expiry path: purchase → expired,
// materialized rows flipped (registrations/product_purchases → expired,
// participations → cancelled), stock freed, notifications fired.
func TestExpire_PendingExpires(t *testing.T) {
	r, w, notifier := newWorker(t)
	purchase, items := seedExpiring(t, r, models.PurchaseStatusPending)

	err := run(t, w, purchase.ID)
	if err != nil {
		t.Fatalf("Work: %v", err)
	}

	// Purchase expired.
	persisted, err := r.Purchases.GetByID(context.Background(), purchase.ID)
	if err != nil {
		t.Fatalf("load purchase: %v", err)
	}
	if persisted.Status != models.PurchaseStatusExpired {
		t.Fatalf("status = %s, want expired", persisted.Status)
	}

	// Materialized rows flipped.
	reg, err := r.Registrations.GetByID(context.Background(), *items[0].RegistrationID)
	if err != nil {
		t.Fatalf("load registration: %v", err)
	}
	if reg.Status != models.RegistrationStatusExpired {
		t.Fatalf("registration = %s, want expired", reg.Status)
	}
	pp, err := r.Products.GetProductPurchaseByID(context.Background(), *items[1].ProductPurchaseID)
	if err != nil {
		t.Fatalf("load product purchase: %v", err)
	}
	if pp.Status != models.ProductPurchaseStatusExpired {
		t.Fatalf("product purchase = %s, want expired", pp.Status)
	}

	// Stock freed: the variant's availability returns to its base (expired
	// purchases don't count as reserved).
	avail, err := r.Purchases.Availability(context.Background(), purchase.EditionID)
	if err != nil {
		t.Fatalf("availability: %v", err)
	}
	for _, a := range avail {
		if a.ItemID == items[1].ItemID {
			if a.ReservedQuantity != 0 {
				t.Fatalf("variant reserved = %d, want 0 (stock freed)", a.ReservedQuantity)
			}
		}
	}

	if notifier.count() != 2 {
		t.Fatalf("notifications = %d, want 2 (stock + purchase.expired)", notifier.count())
	}
}

// TestExpire_ApprovedIsNoOp pins the backstop contract: a purchase already
// approved (webhook won the race; the job cancel missed) is left alone.
func TestExpire_ApprovedIsNoOp(t *testing.T) {
	r, w, notifier := newWorker(t)
	purchase, items := seedExpiring(t, r, models.PurchaseStatusApproved)

	err := run(t, w, purchase.ID)
	if err != nil {
		t.Fatalf("Work: %v", err)
	}

	persisted, err := r.Purchases.GetByID(context.Background(), purchase.ID)
	if err != nil {
		t.Fatalf("load purchase: %v", err)
	}
	if persisted.Status != models.PurchaseStatusApproved {
		t.Fatalf("status = %s, want approved (untouched)", persisted.Status)
	}
	reg, err := r.Registrations.GetByID(context.Background(), *items[0].RegistrationID)
	if err != nil {
		t.Fatalf("load registration: %v", err)
	}
	if reg.Status != models.RegistrationStatusPending {
		t.Fatalf("registration = %s, want pending (untouched)", reg.Status)
	}
	if notifier.count() != 0 {
		t.Fatalf("notifications = %d, want 0 for a no-op", notifier.count())
	}
}

// TestExpire_DoubleRunIsIdempotent pins double-run safety: the second run
// finds the purchase already expired and does nothing.
func TestExpire_DoubleRunIsIdempotent(t *testing.T) {
	r, w, _ := newWorker(t)
	purchase, _ := seedExpiring(t, r, models.PurchaseStatusPending)

	err := run(t, w, purchase.ID)
	if err != nil {
		t.Fatalf("first Work: %v", err)
	}
	err = run(t, w, purchase.ID)
	if err != nil {
		t.Fatalf("second Work: %v", err)
	}

	persisted, err := r.Purchases.GetByID(context.Background(), purchase.ID)
	if err != nil {
		t.Fatalf("load purchase: %v", err)
	}
	if persisted.Status != models.PurchaseStatusExpired {
		t.Fatalf("status = %s, want expired", persisted.Status)
	}
}
