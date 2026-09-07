package badges_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ovechkin-dm/mockio/mock"

	"lib/database"
	"lib/email"

	"univents/models"
	"univents/ports"
)

// fakeTx satisfies pgx.Tx by embedding the interface — it is only ever stored
// in the context to simulate "this context is inside an open transaction",
// never invoked.
type fakeTx struct{ pgx.Tx }

// TestEmitForConfirmedRegistration_EmailGoroutineMustNotShareTheTx is the
// regression guard for the pgx "conn busy" bug: EmitForConfirmedRegistration
// fires the badge email goroutine while the caller may be inside a
// transaction. If the goroutine's context still carries the tx, its
// editions/events reads would run on the caller's transaction connection and
// race the ongoing transaction (pgconn.connLockError "conn busy" and
// SQLSTATE 08P01 on the prepared-statement cache).
//
// The test simulates the real checkout/webhook shape: ctx carries an open tx,
// and the email goroutine's DB reads must NOT observe it.
func TestEmitForConfirmedRegistration_EmailGoroutineMustNotShareTheTx(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	userID := uuid.New()
	regID := uuid.New()

	reg := &models.Registration{
		ID:             regID,
		EditionID:      editionID,
		AttendeeUserID: &userID,
		AttendeeEmail:  "attendee@example.com",
		AttendeeName:   "Attendee",
		Status:         models.RegistrationStatusConfirmed,
	}

	registrations := mock.Mock[ports.RegistrationRepo]()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	editions := mock.Mock[ports.EditionRepo]()
	events := mock.Mock[ports.EventRepo]()

	// Synchronous pre-goroutine calls must still run *with* the tx — proves the
	// test really simulates an in-transaction emit (not that everything is
	// tx-free).
	regsSawTx := make(chan bool, 1)
	mock.When(registrations.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenAnswer(func(args []any) []any {
			regsSawTx <- hasTx(args[0].(context.Context))
			return []any{reg, nil}
		})

	// The emission upsert happens inside the tx too.
	upsertSawTx := make(chan bool, 1)
	mock.When(emissions.Upsert(mock.AnyContext(), mock.Any[*models.BadgeEmission]())).
		ThenAnswer(func(args []any) []any {
			upsertSawTx <- hasTx(args[0].(context.Context))
			em := &models.BadgeEmission{ // EmailSentAt nil → email goroutine fires
				ID: uuid.New(), EditionID: editionID, UserID: userID,
				Origin: models.BadgeEmissionOriginParticipant, RegistrationID: &regID,
			}
			return []any{em, nil}
		})

	// The email goroutine's edition/event lookups: these are the reads that
	// used to race the caller's transaction.
	emailEditionsSawTx := make(chan bool, 1)
	emailEventsSawTx := make(chan bool, 1)
	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenAnswer(func(args []any) []any {
			emailEditionsSawTx <- hasTx(args[0].(context.Context))
			return []any{&models.Edition{ID: editionID, Name: "Edition"}, nil}
		})
	mock.When(events.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenAnswer(func(args []any) []any {
			emailEventsSawTx <- hasTx(args[0].(context.Context))
			return []any{&models.Event{FullName: "Event"}, nil}
		})

	// SMTP dial to a closed port fails harmlessly and is logged, as in
	// emit_test.go — we only need the goroutine to reach the DB reads.
	emailClient := email.NewClient(email.Config{Host: "127.0.0.1", Port: 1, From: "badges@test.local"})

	ops := newOps(t, nil, emissions, registrations, editions, events, emailClient)

	// Simulate emission from inside a transaction (free checkout / webhook
	// approve tx) — ctx carries the open tx.
	txCtx := context.WithValue(context.Background(), database.TxKeyValue, fakeTx{})
	_, err := ops.EmitForConfirmedRegistration(txCtx, regID)
	if err != nil {
		t.Fatalf("EmitForConfirmedRegistration: %v", err)
	}

	if got := waitBool(t, regsSawTx); !got {
		t.Error("setup broken: registration read should have seen the tx")
	}
	if got := waitBool(t, upsertSawTx); !got {
		t.Error("setup broken: emission upsert should have seen the tx")
	}

	// The regression: the email goroutine must not share the tx connection.
	if got := waitBool(t, emailEditionsSawTx); got {
		t.Error("badge email goroutine's edition read ran with the caller's tx in ctx — pgx 'conn busy' race")
	}
	if got := waitBool(t, emailEventsSawTx); got {
		t.Error("badge email goroutine's event read ran with the caller's tx in ctx — pgx 'conn busy' race")
	}
}

func hasTx(ctx context.Context) bool {
	tx, _ := ctx.Value(database.TxKeyValue).(pgx.Tx)
	return tx != nil
}

func waitBool(t *testing.T, ch chan bool) bool {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for mocked DB call")
		return false
	}
}
