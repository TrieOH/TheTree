package webhooks_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river/rivertype"

	"lib/database"
	"lib/testdb"

	"univents/internal/repos"
	"univents/internal/services/webhooks"
	"univents/internal/sqlc"
	"univents/models"
)

// recordingBadges records EmitForConfirmedRegistration calls (the real
// badge emit needs the full badge_emissions machinery; here we assert it
// was called with the right registration).
type recordingBadges struct {
	mu      sync.Mutex
	emitted []uuid.UUID
}

func (b *recordingBadges) EmitForConfirmedRegistration(_ context.Context, id uuid.UUID) (*models.BadgeEmission, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.emitted = append(b.emitted, id)
	return &models.BadgeEmission{ID: uuid.New()}, nil
}

func (b *recordingBadges) ids() []uuid.UUID {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]uuid.UUID{}, b.emitted...)
}

// recordingRiver records JobCancel calls.
type recordingRiver struct {
	mu        sync.Mutex
	cancelled []int64
}

func (r *recordingRiver) JobCancel(_ context.Context, id int64) (*rivertype.JobRow, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cancelled = append(r.cancelled, id)
	return &rivertype.JobRow{}, nil
}

func (r *recordingRiver) ids() []int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]int64{}, r.cancelled...)
}

type dbFixture struct {
	editionID    uuid.UUID
	ticketID     uuid.UUID
	variantID    uuid.UUID
	occurrenceID uuid.UUID
}

func seedDBFixture(t *testing.T, q *sqlc.Queries) dbFixture {
	t.Helper()
	ctx := context.Background()

	event, err := q.CreateEvent(ctx, sqlc.CreateEventParams{
		OwnerID:  uuid.New(),
		FullName: "DB Seed Event",
		Slug:     "db-seed-" + uuid.NewString()[:8],
		Status:   string(models.EventStatusActive),
	})
	if err != nil {
		t.Fatalf("seed event: %v", err)
	}
	edition, err := q.CreateEdition(ctx, sqlc.CreateEditionParams{
		EventID:     event.ID,
		EditionName: "DB Seed Edition",
		Slug:        "db-seed-ed-" + uuid.NewString()[:8],
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
	return dbFixture{edition.ID, ticket.ID, variant.ID, occurrence.ID}
}

// pendingPurchaseSeed is the set of rows a split-7 checkout would create.
type pendingPurchaseSeed struct {
	purchaseID uuid.UUID
	regID      uuid.UUID
	ppID       uuid.UUID
	partID     uuid.UUID
	intentID   uuid.UUID
}

// seedPendingPurchase creates a pending purchase with one item of each type
// and its materialized rows (D4), exactly as split 7 checkout will — in one
// tx.
func seedPendingPurchase(t *testing.T, r *repos.Repos, fx dbFixture) pendingPurchaseSeed {
	t.Helper()
	ctx := context.Background()
	purchaserID := uuid.New()
	intentID := uuid.New()
	jobID := int64(99) // a fake river expiry job, as split 7 would schedule
	var seed pendingPurchaseSeed
	seed.intentID = intentID

	err := database.RunTx(ctx, func(ctx context.Context) error {
		p, err := r.Purchases.CreatePurchase(ctx, &models.Purchase{
			EditionID:        fx.editionID,
			PurchaserID:      purchaserID,
			Status:           models.PurchaseStatusPending,
			TotalCents:       8000,
			Currency:         "BRL",
			ExpiresAt:        time.Now().Add(10 * time.Minute),
			PayssageIntentID: &intentID,
			RiverJobID:       &jobID,
		})
		if err != nil {
			return err
		}
		seed.purchaseID = p.ID

		reg, err := r.Registrations.Create(ctx, &models.Registration{
			EditionID:      fx.editionID,
			TicketTypeID:   fx.ticketID,
			PurchaserID:    purchaserID,
			AttendeeUserID: purchaserID,
			AttendeeEmail:  "buyer@example.com",
			AttendeeName:   "Buyer",
			Status:         models.RegistrationStatusPending,
		})
		if err != nil {
			return err
		}
		seed.regID = reg.ID

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
		seed.ppID = pp.ID

		part, err := r.Programs.CreateParticipation(ctx, &models.ProgramParticipation{
			EditionID:      fx.editionID,
			OccurrenceID:   fx.occurrenceID,
			RegistrationID: reg.ID,
			Status:         models.ProgramParticipationStatusRegistered,
		})
		if err != nil {
			return err
		}
		seed.partID = part.ID

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
		t.Fatalf("seed purchase: %v", err)
	}
	return seed
}

// TestReceiveDB_ApproveEndToEnd exercises the full approve path against a
// real Postgres: signed delivery → guarded pending→approved → materialized
// rows flipped (registrations confirmed, product_purchases confirmed,
// participations stay registered) → badge emit + expiry-job cancel +
// NOTIFY — all committed in one tx.
func TestReceiveDB_ApproveEndToEnd(t *testing.T) {
	pool := testdb.Postgres(t, "../../../db/migrations")
	q := sqlc.New(pool)
	runner := database.NewPGXTxRunner(pool)
	database.SetDefaultRunner(runner)
	r := repos.New(q)

	fx := seedDBFixture(t, q)
	seed := seedPendingPurchase(t, r, fx)

	badges := &recordingBadges{}
	river := &recordingRiver{}
	notifier := &recordingNotifier{}
	ops := webhooks.NewOperations(
		r.Purchases, r.Registrations, r.Products, r.Programs,
		badges, notifier, river, runner, testSecret,
	)
	ops.SetCardRaceWait(0)

	err := ops.Receive(context.Background(), input(seed.intentID, "payment.succeeded", []byte(`{"event":"ok"}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}

	ctx := context.Background()

	got, err := r.Purchases.GetByID(ctx, seed.purchaseID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Status != models.PurchaseStatusApproved {
		t.Fatalf("purchase status = %s, want approved", got.Status)
	}

	reg, err := r.Registrations.GetByID(ctx, seed.regID)
	if err != nil {
		t.Fatalf("registration GetByID: %v", err)
	}
	if reg.Status != models.RegistrationStatusConfirmed {
		t.Fatalf("registration status = %s, want confirmed", reg.Status)
	}

	var ppStatus string
	err = pool.QueryRow(ctx, "SELECT status FROM product_purchases WHERE id = $1", seed.ppID).Scan(&ppStatus)
	if err != nil {
		t.Fatalf("product_purchases read: %v", err)
	}
	if ppStatus != "confirmed" {
		t.Fatalf("product_purchase status = %s, want confirmed", ppStatus)
	}

	// program_participations stay registered on approval (D4).
	var partStatus string
	err = pool.QueryRow(ctx, "SELECT status FROM program_participations WHERE id = $1", seed.partID).Scan(&partStatus)
	if err != nil {
		t.Fatalf("program_participations read: %v", err)
	}
	if partStatus != "registered" {
		t.Fatalf("participation status = %s, want registered", partStatus)
	}

	emitted := badges.ids()
	if len(emitted) != 1 || emitted[0] != seed.regID {
		t.Fatalf("badge emit = %v, want [%s]", emitted, seed.regID)
	}

	cancelled := river.ids()
	if len(cancelled) != 1 || cancelled[0] != 99 {
		t.Fatalf("river JobCancel = %v, want [99]", cancelled)
	}

	notifies := notifier.decoded()
	if len(notifies) != 2 {
		t.Fatalf("notifications = %d, want 2 (stock + purchase)", len(notifies))
	}
}

// TestReceiveDB_DuplicateDeliveryIsNoOp pins idempotency against real SQL:
// a second delivery finds the purchase already approved, the guarded update
// misses, and no materialized row is touched again.
func TestReceiveDB_DuplicateDeliveryIsNoOp(t *testing.T) {
	pool := testdb.Postgres(t, "../../../db/migrations")
	q := sqlc.New(pool)
	runner := database.NewPGXTxRunner(pool)
	database.SetDefaultRunner(runner)
	r := repos.New(q)

	fx := seedDBFixture(t, q)
	seed := seedPendingPurchase(t, r, fx)

	badges := &recordingBadges{}
	river := &recordingRiver{}
	notifier := &recordingNotifier{}
	ops := webhooks.NewOperations(
		r.Purchases, r.Registrations, r.Products, r.Programs,
		badges, notifier, river, runner, testSecret,
	)
	ops.SetCardRaceWait(0)

	ctx := context.Background()
	for i := range 2 {
		err := ops.Receive(ctx, input(seed.intentID, "payment.succeeded", []byte(`{"event":"ok"}`)))
		if err != nil {
			t.Fatalf("Receive #%d: %v", i+1, err)
		}
	}

	got, err := r.Purchases.GetByID(ctx, seed.purchaseID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Status != models.PurchaseStatusApproved {
		t.Fatalf("purchase status = %s, want approved", got.Status)
	}
	// Badge emitted exactly once (upsert is idempotent on re-delivery).
	if ids := badges.ids(); len(ids) != 1 {
		t.Fatalf("badge emits = %d, want 1", len(ids))
	}
	// The second delivery is a full no-op: no re-notify (the first
	// delivery already published stock + purchase).
	if got := notifier.payloads(); len(got) != 2 {
		t.Fatalf("notifications = %d, want 2 (first delivery only)", len(got))
	}
}
