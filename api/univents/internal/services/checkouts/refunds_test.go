package checkouts_test

import (
	"context"
	"testing"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	"lib/database"
	"lib/testdb"

	"univents/internal/authz"
	"univents/internal/repos"
	"univents/internal/services/checkouts"
	"univents/internal/sqlc"
	"univents/models"

	payssage "sdk/payssage"
)

// newRefundOps builds the checkouts operations for the organizer surface
// (refund plan B3) with a seeded event/edition and a fake payssage client.
func newRefundOps(t *testing.T) (*repos.Repos, *checkouts.Operations, *fakes, fixture) {
	t.Helper()
	pool := testdb.Postgres(t, "../../../db/migrations")
	q := sqlc.New(pool)
	tx := database.NewPGXTxRunner(pool)
	database.SetDefaultRunner(tx)
	r := repos.New(q)

	ps := newFakePayssage(func(uuid.UUID, payssage.CreateIntentRequest) (*payssage.Intent, error) {
		return &payssage.Intent{Status: payssage.IntentStatusSucceeded}, nil
	})
	fs := &fakes{
		badges:   &fakeBadges{},
		notifier: &fakeNotifier{},
		river:    &fakeRiver{},
		tokens:   &fakeTokens{},
		payssage: ps,
	}
	ops := checkouts.NewOperations(
		r.Purchases, r.Editions, r.Events, r.TicketTypes, r.Products, r.Programs, r.Occurrences,
		r.Registrations, r.Products, r.Programs,
		fs.badges, fs.notifier, fs.river, tx,
		nil, ps, walletID, fs.tokens, authz.New(r.Events),
	)
	fx := seedStore(t, r)
	// The owner is a member via AddEventMember (Create does not auto-insert
	// the owner row) — required for the organizer authz check.
	_, err := r.Events.AddEventMember(context.Background(), fx.eventID, fx.ownerID, models.EventMemberRoleOwner)
	if err != nil {
		t.Fatalf("seed owner membership: %v", err)
	}
	return r, ops, fs, fx
}

// TestRefundPurchase_OwnerInitiatesRefund pins the organizer refund: owner
// refunds an approved purchase — payssage gets the call and the purchase
// stays approved (the payment.refunded webhook flips it, single writer).
func TestRefundPurchase_OwnerInitiatesRefund(t *testing.T) {
	r, ops, fs, fx := newRefundOps(t)
	intentID := uuid.New()
	purchase, _ := seedPurchase(t, r, fx.editionID, uuid.New(), models.PurchaseStatusApproved, &intentID)

	got, err := ops.RefundPurchase(context.Background(), fx.ownerID, purchase.ID)
	if err != nil {
		t.Fatalf("RefundPurchase: %v", err)
	}
	if got.Status != models.PurchaseStatusApproved {
		t.Fatalf("purchase status = %s, want approved (webhook flips)", got.Status)
	}
	if fs.payssage.refundCount() != 1 {
		t.Fatalf("payssage refunds = %d, want 1", fs.payssage.refundCount())
	}
	dbPurchase, err := r.Purchases.GetByID(context.Background(), purchase.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if dbPurchase.Status != models.PurchaseStatusApproved {
		t.Fatalf("db purchase status = %s, want approved", dbPurchase.Status)
	}
}

// TestRefundPurchase_Guards pins the status guards: already-refunded and
// non-approved purchases cannot be refunded; a stranger cannot refund at all.
func TestRefundPurchase_Guards(t *testing.T) {
	t.Run("already refunded", func(t *testing.T) {
		r, ops, _, fx := newRefundOps(t)
		intentID := uuid.New()
		purchase, _ := seedPurchase(t, r, fx.editionID, uuid.New(), models.PurchaseStatusRefunded, &intentID)
		_, err := ops.RefundPurchase(context.Background(), fx.ownerID, purchase.ID)
		if !fun.Is(err, fun.CodeConflict) {
			t.Fatalf("expected conflict for refunded purchase, got %v", err)
		}
	})
	t.Run("pending", func(t *testing.T) {
		r, ops, _, fx := newRefundOps(t)
		intentID := uuid.New()
		purchase, _ := seedPurchase(t, r, fx.editionID, uuid.New(), models.PurchaseStatusPending, &intentID)
		_, err := ops.RefundPurchase(context.Background(), fx.ownerID, purchase.ID)
		if !fun.Is(err, fun.CodeConflict) {
			t.Fatalf("expected conflict for pending purchase, got %v", err)
		}
	})
	t.Run("stranger forbidden", func(t *testing.T) {
		r, ops, fs, fx := newRefundOps(t)
		intentID := uuid.New()
		purchase, _ := seedPurchase(t, r, fx.editionID, uuid.New(), models.PurchaseStatusApproved, &intentID)
		_, err := ops.RefundPurchase(context.Background(), uuid.New(), purchase.ID)
		if err == nil {
			t.Fatal("expected error for non-member")
		}
		if fs.payssage.refundCount() != 0 {
			t.Fatalf("payssage refunds = %d, want 0", fs.payssage.refundCount())
		}
	})
}

// TestListEditionPurchases_OwnerSeesPurchases pins the organizer orders read:
// the owner sees the edition's purchases with their items; a stranger gets a
// forbidden/not-found error and sees nothing.
func TestListEditionPurchases_OwnerSeesPurchases(t *testing.T) {
	r, ops, _, fx := newRefundOps(t)
	intentID := uuid.New()
	purchase, items := seedPurchase(t, r, fx.editionID, uuid.New(), models.PurchaseStatusApproved, &intentID)

	rows, err := ops.ListEditionPurchases(context.Background(), fx.ownerID, fx.editionID)
	if err != nil {
		t.Fatalf("ListEditionPurchases: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("purchases = %d, want 1", len(rows))
	}
	if rows[0].Purchase.ID != purchase.ID {
		t.Fatalf("purchase id = %s, want %s", rows[0].Purchase.ID, purchase.ID)
	}
	if len(rows[0].Items) != len(items) {
		t.Fatalf("items = %d, want %d", len(rows[0].Items), len(items))
	}
}

func TestListEditionPurchases_StrangerForbidden(t *testing.T) {
	r, ops, _, fx := newRefundOps(t)
	intentID := uuid.New()
	seedPurchase(t, r, fx.editionID, uuid.New(), models.PurchaseStatusApproved, &intentID)

	_, err := ops.ListEditionPurchases(context.Background(), uuid.New(), fx.editionID)
	if err == nil {
		t.Fatal("expected error for non-member")
	}
}
