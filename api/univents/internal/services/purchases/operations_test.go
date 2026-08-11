package purchases_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"lib/database"
	"lib/testdb"

	"univents/internal/repos"
	"univents/internal/services/purchases"
	"univents/internal/sqlc"
	"univents/models"
)

// seedEdition creates an event + edition the purchases can reference.
func seedEdition(t *testing.T, q *sqlc.Queries) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	event, err := q.CreateEvent(ctx, sqlc.CreateEventParams{
		OwnerID:  uuid.New(),
		FullName: "List Test Event",
		Slug:     "list-" + uuid.NewString()[:8],
		Status:   string(models.EventStatusActive),
	})
	if err != nil {
		t.Fatalf("seed event: %v", err)
	}
	edition, err := q.CreateEdition(ctx, sqlc.CreateEditionParams{
		EventID:     event.ID,
		EditionName: "List Test Edition",
		Slug:        "list-ed-" + uuid.NewString()[:8],
		StartsAt:    time.Now().Add(-time.Hour),
		EndsAt:      time.Now().Add(24 * time.Hour),
		CreatedBy:   event.OwnerID,
	})
	if err != nil {
		t.Fatalf("seed edition: %v", err)
	}
	return edition.ID
}

// seedPurchase creates one purchase (+ one ticket item) for the user. The
// optional sleep delays the insert so created_at (DEFAULT now(), microsecond
// resolution) is deterministically ordered across calls.
func seedPurchase(t *testing.T, r *repos.Repos, editionID, purchaserID uuid.UUID, status models.PurchaseStatus, sleep time.Duration) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	time.Sleep(sleep)

	purchase, err := r.Purchases.CreatePurchase(ctx, &models.Purchase{
		EditionID:   editionID,
		PurchaserID: purchaserID,
		Status:      status,
		TotalCents:  1500,
		Currency:    "BRL",
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	})
	if err != nil {
		t.Fatalf("seed purchase: %v", err)
	}
	created, err := r.Purchases.CreatePurchaseItem(ctx, &models.PurchaseItem{
		PurchaseID:     purchase.ID,
		ItemType:       models.PurchaseItemTypeTicket,
		ItemID:         uuid.New(),
		Quantity:       1,
		UnitPriceCents: 1500,
	})
	if err != nil {
		t.Fatalf("seed purchase item: %v", err)
	}
	_ = created
	return purchase.ID
}

func newOps(t *testing.T) (*repos.Repos, *sqlc.Queries, *purchases.Operations) {
	t.Helper()
	pool := testdb.Postgres(t, "../../../db/migrations")
	q := sqlc.New(pool)
	database.SetDefaultRunner(database.NewPGXTxRunner(pool))
	r := repos.New(q)
	return r, q, purchases.NewOperations(r.Purchases)
}

// TestListForUser_OnlyOwnPurchasesNewestFirst pins the core contract: the
// caller sees exactly their own purchases, newest first, each with its
// items.
func TestListForUser_OnlyOwnPurchasesNewestFirst(t *testing.T) {
	r, q, ops := newOps(t)
	editionID := seedEdition(t, q)

	owner := uuid.New()
	other := uuid.New()
	// owner's older purchase first, then a 20ms gap, then the newer one —
	// the list must come back newest first.
	older := seedPurchase(t, r, editionID, owner, models.PurchaseStatusPending, 0)
	_ = seedPurchase(t, r, editionID, other, models.PurchaseStatusPending, 10*time.Millisecond) // someone else's
	newer := seedPurchase(t, r, editionID, owner, models.PurchaseStatusApproved, 20*time.Millisecond)

	details, err := ops.ListForUser(context.Background(), owner)
	if err != nil {
		t.Fatalf("ListForUser: %v", err)
	}

	if len(details) != 2 {
		t.Fatalf("purchases = %d, want 2 (only the owner's)", len(details))
	}
	if details[0].Purchase.ID != newer {
		t.Fatalf("first purchase = %s, want newest %s", details[0].Purchase.ID, newer)
	}
	if details[1].Purchase.ID != older {
		t.Fatalf("second purchase = %s, want older %s", details[1].Purchase.ID, older)
	}
	if details[1].Purchase.Status != models.PurchaseStatusPending {
		t.Fatalf("older purchase status = %s, want pending", details[1].Purchase.Status)
	}

	// Items ride along with each purchase.
	for _, d := range details {
		if len(d.Items) != 1 {
			t.Fatalf("purchase %s items = %d, want 1", d.Purchase.ID, len(d.Items))
		}
		if d.Items[0].ItemType != models.PurchaseItemTypeTicket {
			t.Fatalf("purchase %s item type = %s, want ticket", d.Purchase.ID, d.Items[0].ItemType)
		}
	}
}

// TestListForUser_Empty pins the no-history case: a fresh user gets an
// empty (non-nil) list — the handler marshals it as `[]`, never `null`.
func TestListForUser_Empty(t *testing.T) {
	_, _, ops := newOps(t)

	details, err := ops.ListForUser(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("ListForUser: %v", err)
	}
	if details == nil || len(details) != 0 {
		t.Fatalf("purchases = %#v, want empty non-nil list", details)
	}
}
