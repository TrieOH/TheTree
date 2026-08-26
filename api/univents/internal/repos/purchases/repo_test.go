package purchases_test

import (
	"context"
	"testing"
	"time"

	"lib/database"
	"lib/testdb"
	"univents/internal/repos"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
)

// fixture is the minimum store graph: event → edition → ticket_type +
// product/variant + program/occurrence.
type fixture struct {
	editionID    uuid.UUID
	ticketID     uuid.UUID
	variantID    uuid.UUID
	occurrenceID uuid.UUID
}

func seedFixture(t *testing.T, q *sqlc.Queries) fixture {
	t.Helper()
	ctx := context.Background()

	event, err := q.CreateEvent(ctx, sqlc.CreateEventParams{
		OwnerID:  uuid.New(),
		FullName: "Seed Event",
		Slug:     "seed-" + uuid.NewString()[:8],
		Status:   string(models.EventStatusActive),
	})
	if err != nil {
		t.Fatalf("seed event: %v", err)
	}

	edition, err := q.CreateEdition(ctx, sqlc.CreateEditionParams{
		EventID:     event.ID,
		EditionName: "Seed Edition",
		Slug:        "seed-ed-" + uuid.NewString()[:8],
		StartsAt:    time.Now().Add(-time.Hour),
		EndsAt:      time.Now().Add(24 * time.Hour),
		CreatedBy:   event.OwnerID,
	})
	if err != nil {
		t.Fatalf("seed edition: %v", err)
	}

	ticket, err := q.CreateTicketType(ctx, sqlc.CreateTicketTypeParams{
		EditionID:   edition.ID,
		Name:        "Standard",
		AccessLevel: 0,
		Price:       1000,
		MaxQuantity: new(int(10)),
	})
	if err != nil {
		t.Fatalf("seed ticket: %v", err)
	}

	product, err := q.CreateProduct(ctx, sqlc.CreateProductParams{
		EditionID:  edition.ID,
		VendorCode: "P-" + uuid.NewString()[:8],
	})
	if err != nil {
		t.Fatalf("seed product: %v", err)
	}
	variant, err := q.CreateProductVariant(ctx, sqlc.CreateProductVariantParams{
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

	program, err := q.CreateProgram(ctx, sqlc.CreateProgramParams{
		EditionID: edition.ID,
		Kind:      string(models.ProgramKindActivity),
		Name:      "Workshop",
		Price:     2000,
	})
	if err != nil {
		t.Fatalf("seed program: %v", err)
	}
	occurrence, err := q.CreateProgramOccurrence(ctx, sqlc.CreateProgramOccurrenceParams{
		ProgramID:   program.ID,
		EditionID:   edition.ID,
		StartsAt:    time.Now().Add(2 * time.Hour),
		EndsAt:      time.Now().Add(3 * time.Hour),
		MaxCapacity: new(int(3)),
	})
	if err != nil {
		t.Fatalf("seed occurrence: %v", err)
	}

	return fixture{edition.ID, ticket.ID, variant.ID, occurrence.ID}
}

func newRepo(t *testing.T) (*repos.Repos, *sqlc.Queries) {
	t.Helper()
	pool := testdb.Postgres(t, "../../../db/migrations")
	q := sqlc.New(pool)
	database.SetDefaultRunner(database.NewPGXTxRunner(pool))
	return repos.New(q), q
}

func purchase(id uuid.UUID, editionID uuid.UUID, status models.PurchaseStatus) *models.Purchase {
	return &models.Purchase{
		EditionID:   editionID,
		PurchaserID: id,
		Status:      status,
		TotalCents:  8000,
		Currency:    "BRL",
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}
}

// availabilityByID indexes an availability listing by item id for assertions.
func availabilityByID(t *testing.T, r *repos.Repos, editionID uuid.UUID) map[uuid.UUID]models.ItemAvailability {
	t.Helper()
	avail, err := r.Purchases.Availability(context.Background(), editionID)
	if err != nil {
		t.Fatalf("availability: %v", err)
	}
	byID := make(map[uuid.UUID]models.ItemAvailability, len(avail))
	for _, a := range avail {
		byID[a.ItemID] = a
	}
	return byID
}

// TestCreatePurchaseWithItemsInOneTx pins the D4 contract: purchase + items +
// materialized rows (registration / product purchase / participation) commit
// in one tx, and the materialized ids are stored back on the items.
func TestCreatePurchaseWithItemsInOneTx(t *testing.T) {
	r, q := newRepo(t)
	fx := seedFixture(t, q)
	purchaserID := uuid.New()
	ctx := context.Background()

	var purchaseID uuid.UUID
	err := database.RunTx(ctx, func(ctx context.Context) error {
		p, err := r.Purchases.CreatePurchase(ctx, purchase(purchaserID, fx.editionID, models.PurchaseStatusPending))
		if err != nil {
			return err
		}
		purchaseID = p.ID

		reg, err := r.Registrations.Create(ctx, &models.Registration{
			EditionID:      fx.editionID,
			TicketTypeID:   fx.ticketID,
			PurchaserID:    purchaserID,
			AttendeeUserID: &purchaserID,
			AttendeeEmail:  "buyer@example.com",
			AttendeeName:   "Buyer",
			Status:         models.RegistrationStatusPending,
		})
		if err != nil {
			return err
		}

		pp, err := r.Products.CreateProductPurchase(ctx, &models.ProductPurchase{
			EditionID:      fx.editionID,
			VariantID:      fx.variantID,
			PurchaserID:    purchaserID,
			RegistrationID: &reg.ID,
			Quantity:       1,
			Status:         models.ProductPurchaseStatusPending,
		})
		if err != nil {
			return err
		}

		part, err := r.Programs.CreateParticipation(ctx, &models.ProgramParticipation{
			EditionID:      fx.editionID,
			OccurrenceID:   fx.occurrenceID,
			RegistrationID: reg.ID,
			Status:         models.ProgramParticipationStatusRegistered,
		})
		if err != nil {
			return err
		}

		for _, item := range []*models.PurchaseItem{
			{ItemType: models.PurchaseItemTypeTicket, ItemID: fx.ticketID, Quantity: 1, UnitPriceCents: 1000, RegistrationID: &reg.ID},
			{ItemType: models.PurchaseItemTypeProduct, ItemID: fx.variantID, Quantity: 1, UnitPriceCents: 5000, ProductPurchaseID: &pp.ID},
			{ItemType: models.PurchaseItemTypeProgramOccurrence, ItemID: fx.occurrenceID, Quantity: 1, UnitPriceCents: 2000, ParticipationID: &part.ID},
		} {
			item.PurchaseID = p.ID
			_, err := r.Purchases.CreatePurchaseItem(ctx, item)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("tx: %v", err)
	}

	got, err := r.Purchases.GetByID(ctx, purchaseID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Status != models.PurchaseStatusPending || got.TotalCents != 8000 {
		t.Fatalf("purchase = %+v", got)
	}

	items, err := r.Purchases.ListItemsByPurchase(ctx, purchaseID)
	if err != nil {
		t.Fatalf("ListItems: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("items = %d, want 3", len(items))
	}
	for _, item := range items {
		switch item.ItemType {
		case models.PurchaseItemTypeTicket:
			if item.RegistrationID == nil {
				t.Fatal("ticket item missing registration link")
			}
		case models.PurchaseItemTypeProduct:
			if item.ProductPurchaseID == nil {
				t.Fatal("product item missing product purchase link")
			}
		case models.PurchaseItemTypeProgramOccurrence:
			if item.ParticipationID == nil {
				t.Fatal("program item missing participation link")
			}
		}
	}
}

func TestGetByIntentID(t *testing.T) {
	r, q := newRepo(t)
	fx := seedFixture(t, q)
	intentID := uuid.New()
	ctx := context.Background()

	p, err := r.Purchases.CreatePurchase(ctx, &models.Purchase{
		EditionID:        fx.editionID,
		PurchaserID:      uuid.New(),
		Status:           models.PurchaseStatusPending,
		TotalCents:       1000,
		Currency:         "BRL",
		ExpiresAt:        time.Now().Add(10 * time.Minute),
		PayssageIntentID: &intentID,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := r.Purchases.GetByIntentID(ctx, intentID)
	if err != nil {
		t.Fatalf("GetByIntentID: %v", err)
	}
	if got.ID != p.ID {
		t.Fatalf("id = %s, want %s", got.ID, p.ID)
	}
}

func TestOwnerScoping(t *testing.T) {
	r, q := newRepo(t)
	fx := seedFixture(t, q)
	owner, other := uuid.New(), uuid.New()
	ctx := context.Background()

	p, err := r.Purchases.CreatePurchase(ctx, purchase(owner, fx.editionID, models.PurchaseStatusPending))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	_, err = r.Purchases.GetByIDForOwner(ctx, p.ID, owner)
	if err != nil {
		t.Fatalf("owner must see their purchase: %v", err)
	}
	_, err = r.Purchases.GetByIDForOwner(ctx, p.ID, other)
	if err == nil {
		t.Fatal("other user must not see the purchase")
	}

	// GetByID (internal) is owner-agnostic.
	_, err = r.Purchases.GetByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
}

// TestAvailabilityMath pins the split-3 stock semantics:
// available = base - reserved, nil base = unlimited. Reserved counts
// purchase_items of pending AND approved purchases; expired frees stock.
func TestAvailabilityMath(t *testing.T) {
	r, q := newRepo(t)
	fx := seedFixture(t, q)
	buyer := uuid.New()
	ctx := context.Background()

	// No reservations yet: base is untouched.
	byID := availabilityByID(t, r, fx.editionID)
	if a := byID[fx.ticketID]; a.BaseQuantity == nil || *a.BaseQuantity != 10 || a.ReservedQuantity != 0 {
		t.Fatalf("ticket availability = %+v", a)
	}
	if a := byID[fx.variantID]; a.BaseQuantity == nil || *a.BaseQuantity != 5 || a.ReservedQuantity != 0 {
		t.Fatalf("variant availability = %+v", a)
	}
	if a := byID[fx.occurrenceID]; a.BaseQuantity == nil || *a.BaseQuantity != 3 || a.ReservedQuantity != 0 {
		t.Fatalf("occurrence availability = %+v", a)
	}

	// Reserve 2 tickets + 1 variant via a pending purchase.
	var purchaseID uuid.UUID
	err := database.RunTx(ctx, func(ctx context.Context) error {
		p, err := r.Purchases.CreatePurchase(ctx, purchase(buyer, fx.editionID, models.PurchaseStatusPending))
		if err != nil {
			return err
		}
		purchaseID = p.ID
		for _, item := range []*models.PurchaseItem{
			{ItemType: models.PurchaseItemTypeTicket, ItemID: fx.ticketID, Quantity: 2, UnitPriceCents: 1000},
			{ItemType: models.PurchaseItemTypeProduct, ItemID: fx.variantID, Quantity: 1, UnitPriceCents: 5000},
		} {
			item.PurchaseID = p.ID
			_, err := r.Purchases.CreatePurchaseItem(ctx, item)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("reserve tx: %v", err)
	}

	byID = availabilityByID(t, r, fx.editionID)
	if a := byID[fx.ticketID]; a.ReservedQuantity != 2 {
		t.Fatalf("ticket reserved = %d, want 2", a.ReservedQuantity)
	}
	if a := byID[fx.variantID]; a.ReservedQuantity != 1 {
		t.Fatalf("variant reserved = %d, want 1", a.ReservedQuantity)
	}

	// Expired purchases release their reservation.
	_, err = r.Purchases.UpdateStatus(ctx, purchaseID, models.PurchaseStatusExpired, nil)
	if err != nil {
		t.Fatalf("expire: %v", err)
	}
	byID = availabilityByID(t, r, fx.editionID)
	if a := byID[fx.ticketID]; a.ReservedQuantity != 0 {
		t.Fatalf("ticket reserved after expire = %d, want 0", a.ReservedQuantity)
	}
}

// TestAvailabilityUnlimitedPins nil base = unlimited (never sells out).
func TestAvailabilityUnlimited(t *testing.T) {
	r, q := newRepo(t)
	fx := seedFixture(t, q)
	ctx := context.Background()

	var unlimitedID uuid.UUID
	err := database.RunTx(ctx, func(ctx context.Context) error {
		tt, err := q.CreateTicketType(ctx, sqlc.CreateTicketTypeParams{
			EditionID:   fx.editionID,
			Name:        "Unlimited",
			AccessLevel: 0,
			Price:       0,
		})
		if err != nil {
			return err
		}
		unlimitedID = tt.ID
		return nil
	})
	if err != nil {
		t.Fatalf("seed unlimited ticket: %v", err)
	}

	for itemID, a := range availabilityByID(t, r, fx.editionID) {
		if itemID == unlimitedID && a.BaseQuantity != nil {
			t.Fatalf("unlimited ticket base = %v, want nil", a.BaseQuantity)
		}
	}
}

// TestPartialUniqueBlocksSecondPendingPurchase pins the checkout 409: one
// pending purchase per (purchaser, edition).
func TestPartialUniqueBlocksSecondPendingPurchase(t *testing.T) {
	r, q := newRepo(t)
	fx := seedFixture(t, q)
	buyer := uuid.New()
	ctx := context.Background()

	_, err := r.Purchases.CreatePurchase(ctx, purchase(buyer, fx.editionID, models.PurchaseStatusPending))
	if err != nil {
		t.Fatalf("first purchase: %v", err)
	}
	_, err = r.Purchases.CreatePurchase(ctx, purchase(buyer, fx.editionID, models.PurchaseStatusPending))
	if err == nil {
		t.Fatal("second pending purchase must fail (partial unique index)")
	}

	// A different buyer may have their own pending purchase.
	_, err = r.Purchases.CreatePurchase(ctx, purchase(uuid.New(), fx.editionID, models.PurchaseStatusPending))
	if err != nil {
		t.Fatalf("other buyer: %v", err)
	}
}

func TestListByPurchaser(t *testing.T) {
	r, q := newRepo(t)
	fx := seedFixture(t, q)
	buyer := uuid.New()
	ctx := context.Background()

	// Distinct statuses: the partial unique index allows one pending per
	// (purchaser, edition), but list is not status-filtered.
	for i, status := range []models.PurchaseStatus{
		models.PurchaseStatusPending,
		models.PurchaseStatusApproved,
		models.PurchaseStatusExpired,
	} {
		_, err := r.Purchases.CreatePurchase(ctx, purchase(buyer, fx.editionID, status))
		if err != nil {
			t.Fatalf("create %d: %v", i, err)
		}
	}

	list, err := r.Purchases.ListByPurchaser(ctx, buyer)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("list = %d, want 3", len(list))
	}
}
