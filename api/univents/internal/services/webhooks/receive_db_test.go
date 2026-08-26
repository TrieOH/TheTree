package webhooks_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"

	"lib/database"
	"lib/testdb"

	"univents/internal/repos"
	"univents/internal/services/checkouts/jobs"
	"univents/internal/services/webhooks"
	"univents/internal/sqlc"
	"univents/models"
)

// recordingBadges records EmitForConfirmedRegistration and
// RevokeForRegistration calls (the real badge machinery needs the full
// badge_emissions DB; here we assert the calls happened with the right
// registration ids).
type recordingBadges struct {
	mu      sync.Mutex
	emitted []uuid.UUID
	revoked []uuid.UUID
}

func (b *recordingBadges) EmitForConfirmedRegistration(_ context.Context, id uuid.UUID) (*models.BadgeEmission, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.emitted = append(b.emitted, id)
	return &models.BadgeEmission{ID: uuid.New()}, nil
}

func (b *recordingBadges) RevokeForRegistration(_ context.Context, id uuid.UUID, _ string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.revoked = append(b.revoked, id)
	return nil
}

func (b *recordingBadges) ids() []uuid.UUID {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]uuid.UUID{}, b.emitted...)
}

func (b *recordingBadges) revokedIDs() []uuid.UUID {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]uuid.UUID{}, b.revoked...)
}

// recordingRiver records JobCancel calls and captures InsertTx args (the
// gifted-ticket email enqueue).
type recordingRiver struct {
	mu        sync.Mutex
	cancelled []int64
	inserted  []rivertype.JobArgs
}

func (r *recordingRiver) JobCancel(_ context.Context, id int64) (*rivertype.JobRow, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cancelled = append(r.cancelled, id)
	return &rivertype.JobRow{}, nil
}

func (r *recordingRiver) InsertTx(_ context.Context, _ pgx.Tx, args river.JobArgs, _ *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.inserted = append(r.inserted, args)
	return &rivertype.JobInsertResult{}, nil
}

func (r *recordingRiver) ids() []int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]int64{}, r.cancelled...)
}

func (r *recordingRiver) insertedArgs() []rivertype.JobArgs {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]rivertype.JobArgs{}, r.inserted...)
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
			AttendeeUserID: &purchaserID,
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
// notifications.
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

	// Account holders never get the gift email — no gifts.send_email job.
	if got := river.insertedArgs(); len(got) != 0 {
		t.Fatalf("river inserts = %v, want none (attendee has an account)", got)
	}
}

// TestReceiveDB_ApproveGiftEmailOnly pins the gifted-ticket email trigger:
// approving a purchase whose ticket attendee has no account yet (email-only
// gift) confirms the registration, defers the badge, and enqueues
// gifts.send_email in the same tx.
func TestReceiveDB_ApproveGiftEmailOnly(t *testing.T) {
	pool := testdb.Postgres(t, "../../../db/migrations")
	q := sqlc.New(pool)
	runner := database.NewPGXTxRunner(pool)
	database.SetDefaultRunner(runner)
	r := repos.New(q)

	fx := seedDBFixture(t, q)

	// A pending purchase whose only item is an email-only gift (the
	// split-7 checkout shape for an accountless recipient).
	ctx := context.Background()
	purchaserID := uuid.New()
	intentID := uuid.New()
	var regID uuid.UUID
	err := database.RunTx(ctx, func(ctx context.Context) error {
		p, err := r.Purchases.CreatePurchase(ctx, &models.Purchase{
			EditionID:        fx.editionID,
			PurchaserID:      purchaserID,
			Status:           models.PurchaseStatusPending,
			TotalCents:       1000,
			Currency:         "BRL",
			ExpiresAt:        time.Now().Add(10 * time.Minute),
			PayssageIntentID: &intentID,
		})
		if err != nil {
			return err
		}
		reg, err := r.Registrations.Create(ctx, &models.Registration{
			EditionID:      fx.editionID,
			TicketTypeID:   fx.ticketID,
			PurchaserID:    purchaserID,
			AttendeeUserID: nil, // accountless recipient
			AttendeeEmail:  "friend@example.com",
			AttendeeName:   "John Doe",
			Status:         models.RegistrationStatusPending,
		})
		if err != nil {
			return err
		}
		regID = reg.ID
		_, err = r.Purchases.CreatePurchaseItem(ctx, &models.PurchaseItem{
			PurchaseID:     p.ID,
			ItemType:       models.PurchaseItemTypeTicket,
			ItemID:         fx.ticketID,
			Quantity:       1,
			UnitPriceCents: 1000,
			RegistrationID: &reg.ID,
		})
		return err
	})
	if err != nil {
		t.Fatalf("seed gift purchase: %v", err)
	}

	badges := &recordingBadges{}
	river := &recordingRiver{}
	ops := webhooks.NewOperations(
		r.Purchases, r.Registrations, r.Products, r.Programs,
		badges, &recordingNotifier{}, river, runner, testSecret,
	)
	ops.SetCardRaceWait(0)

	err = ops.Receive(ctx, input(intentID, "payment.succeeded", []byte(`{"event":"ok"}`)))
	if err != nil {
		t.Fatalf("Receive: %v", err)
	}

	reg, err := r.Registrations.GetByID(ctx, regID)
	if err != nil {
		t.Fatalf("registration GetByID: %v", err)
	}
	if reg.Status != models.RegistrationStatusConfirmed {
		t.Fatalf("registration status = %s, want confirmed", reg.Status)
	}
	if reg.AttendeeUserID != nil {
		t.Fatalf("attendee_user_id = %v, want nil (still unclaimed)", *reg.AttendeeUserID)
	}

	// Badge emission is called but deferred by the real badges service for
	// accountless registrations (no profile to attach it to) — that skip is
	// covered by the badges unit test; the emit happens for real when the
	// recipient claims and the my-ticket read re-runs it.
	if emitted := badges.ids(); len(emitted) != 1 || emitted[0] != regID {
		t.Fatalf("badge emit = %v, want [%s] (deferred inside the badges service)", emitted, regID)
	}

	// The gift email job is enqueued atomically with the confirmation.
	inserted := river.insertedArgs()
	if len(inserted) != 1 {
		t.Fatalf("river inserts = %d, want 1 (gifts.send_email)", len(inserted))
	}
	args, ok := inserted[0].(jobs.SendGiftEmailArgs)
	if !ok {
		t.Fatalf("inserted args = %T, want jobs.SendGiftEmailArgs", inserted[0])
	}
	if args.RegistrationID != regID {
		t.Fatalf("gift email registration = %s, want %s", args.RegistrationID, regID)
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

// TestReceiveDB_RefundEndToEnd drives the full refund path (refund plan
// slice 3): approve via payment.succeeded, then flip via payment.refunded —
// purchase → refunded, materialized rows → cancelled, badge revoked, stock
// NOTIFYed (no purchase event).
func TestReceiveDB_RefundEndToEnd(t *testing.T) {
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

	err := ops.Receive(ctx, input(seed.intentID, "payment.succeeded", []byte(`{"event":"ok"}`)))
	if err != nil {
		t.Fatalf("approve Receive: %v", err)
	}
	err = ops.Receive(ctx, input(seed.intentID, "payment.refunded", []byte(`{"event":"refund"}`)))
	if err != nil {
		t.Fatalf("refund Receive: %v", err)
	}

	got, err := r.Purchases.GetByID(ctx, seed.purchaseID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Status != models.PurchaseStatusRefunded {
		t.Fatalf("purchase status = %s, want refunded", got.Status)
	}

	reg, err := r.Registrations.GetByID(ctx, seed.regID)
	if err != nil {
		t.Fatalf("registration GetByID: %v", err)
	}
	if reg.Status != models.RegistrationStatusCancelled {
		t.Fatalf("registration status = %s, want cancelled", reg.Status)
	}

	var ppStatus string
	err = pool.QueryRow(ctx, "SELECT status FROM product_purchases WHERE id = $1", seed.ppID).Scan(&ppStatus)
	if err != nil {
		t.Fatalf("product_purchases read: %v", err)
	}
	if ppStatus != "cancelled" {
		t.Fatalf("product_purchase status = %s, want cancelled", ppStatus)
	}

	var partStatus string
	err = pool.QueryRow(ctx, "SELECT status FROM program_participations WHERE id = $1", seed.partID).Scan(&partStatus)
	if err != nil {
		t.Fatalf("program_participations read: %v", err)
	}
	if partStatus != "cancelled" {
		t.Fatalf("participation status = %s, want cancelled", partStatus)
	}

	// Badge emitted on approve, revoked on refund.
	if ids := badges.ids(); len(ids) != 1 || ids[0] != seed.regID {
		t.Fatalf("badge emits = %v, want [%s]", ids, seed.regID)
	}
	if ids := badges.revokedIDs(); len(ids) != 1 || ids[0] != seed.regID {
		t.Fatalf("badge revokes = %v, want [%s]", ids, seed.regID)
	}

	// Approve published stock + purchase (2); refund publishes stock only (1).
	if got := notifier.payloads(); len(got) != 3 {
		t.Fatalf("notifications = %d, want 3 (approve stock+purchase, refund stock)", len(got))
	}
	last := notifier.payloads()[2]
	if !strings.Contains(last, `"kind":"stock"`) {
		t.Fatalf("refund notification = %s, want stock-only (no purchase event)", last)
	}
}
