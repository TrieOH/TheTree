package jobs

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

// fakeTx satisfies pgx.Tx by embedding the interface — only ever stored in a
// context to simulate being inside an open transaction, never invoked.
type fakeTx struct{ pgx.Tx }

// TestEmitCerts_EmailGoroutineMustNotShareTheTx is the regression guard for
// the pgx "conn busy" bug in the cert-emission goroutine: emitCerts fires
// sendCertEmail asynchronously with the caller's context. When the emission
// runs inside a transaction, the goroutine's editions/events reads must fall
// back to the pool instead of racing the caller's transaction connection.
func TestEmitCerts_EmailGoroutineMustNotShareTheTx(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	attendee := models.CertEligibleAttendee{
		UserID:         uuid.New(),
		RegistrationID: uuid.New(),
		AttendeeEmail:  "attendee@example.com",
		AttendeeName:   "Attendee",
	}

	certs := mock.Mock[ports.CertificationRepo]()
	editions := mock.Mock[ports.EditionRepo]()
	events := mock.Mock[ports.EventRepo]()

	// Certify is synchronous inside emitCerts — it runs with the caller's ctx
	// (the tx is legitimately present there).
	certifySawTx := make(chan bool, 1)
	cert := &models.Certification{ID: uuid.New(), EditionID: editionID, VerificationHash: "hash"}
	mock.When(certs.Certify(mock.AnyContext(), mock.Any[models.CertifyInput]())).
		ThenAnswer(func(args []any) []any {
			certifySawTx <- hasTxValue(args[0].(context.Context))
			return []any{cert, nil}
		})

	// The email goroutine's reads: these must NOT see the caller's tx.
	emailEditionsSawTx := make(chan bool, 1)
	emailEventsSawTx := make(chan bool, 1)
	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenAnswer(func(args []any) []any {
			emailEditionsSawTx <- hasTxValue(args[0].(context.Context))
			return []any{&models.Edition{ID: editionID, Name: "Edition"}, nil}
		})
	mock.When(events.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenAnswer(func(args []any) []any {
			emailEventsSawTx <- hasTxValue(args[0].(context.Context))
			return []any{&models.Event{FullName: "Event"}, nil}
		})

	// SMTP dial to a closed port fails harmlessly and is logged; the goroutine
	// only needs to reach the DB reads for the assertion.
	emailClient := email.NewClient(email.Config{Host: "127.0.0.1", Port: 1, From: "certs@test.local"})

	w := NewGrantCertsWorker(certs, editions, events, emailClient)

	txCtx := context.WithValue(context.Background(), database.TxKeyValue, fakeTx{})
	errs := w.emitCerts(txCtx, editionID, nil, []models.CertEligibleAttendee{attendee}, nil)
	if len(errs) != 0 {
		t.Fatalf("emitCerts returned errors: %v", errs)
	}

	if got := waitCtxBool(t, certifySawTx); !got {
		t.Error("setup broken: Certify should have seen the tx")
	}

	// The regression: the async email must not share the caller's tx conn.
	if got := waitCtxBool(t, emailEditionsSawTx); got {
		t.Error("cert email goroutine's edition read ran with the caller's tx in ctx — pgx 'conn busy' race")
	}
	if got := waitCtxBool(t, emailEventsSawTx); got {
		t.Error("cert email goroutine's event read ran with the caller's tx in ctx — pgx 'conn busy' race")
	}
}

func hasTxValue(ctx context.Context) bool {
	tx, _ := ctx.Value(database.TxKeyValue).(pgx.Tx)
	return tx != nil
}

func waitCtxBool(t *testing.T, ch chan bool) bool {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for mocked DB call")
		return false
	}
}
