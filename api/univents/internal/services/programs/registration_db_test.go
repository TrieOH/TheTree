package programs_test

import (
	"context"
	"testing"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"lib/database"
	"lib/testdb"

	"univents/internal/authz"
	"univents/internal/repos"
	"univents/internal/services/programs"
	"univents/internal/sqlc"
	"univents/models"
)

// fixture is the minimum activity graph: event → edition → ticket type →
// confirmed registration → program → occurrence with capacity 2.
type fixture struct {
	eventID      uuid.UUID
	ownerID      uuid.UUID
	editionID    uuid.UUID
	ticketID     uuid.UUID
	regA         uuid.UUID
	attendeeA    uuid.UUID
	regB         uuid.UUID
	attendeeB    uuid.UUID
	regC         uuid.UUID
	attendeeC    uuid.UUID
	programID    uuid.UUID
	occurrenceID uuid.UUID
	staffID      uuid.UUID

	// checkpoint fixture (TestCheckpointCheckIn_*)
	checkpointProgramID    uuid.UUID
	checkpointOccurrenceID uuid.UUID
}

// seedActivity creates the fixture through the real repos (disposable
// Postgres with the real migrations). Each attendee holds their own
// confirmed registration; staffID is an event staff member.
func seedActivity(t *testing.T, r *repos.Repos) fixture {
	t.Helper()
	ctx := context.Background()
	fx := fixture{
		ownerID:   uuid.New(),
		attendeeA: uuid.New(),
		attendeeB: uuid.New(),
		attendeeC: uuid.New(),
		staffID:   uuid.New(),
	}

	event, err := r.Events.Create(ctx, &models.Event{
		OwnerID: fx.ownerID, FullName: "Activity Test Event",
		Slug: "activity-test-" + uuid.NewString()[:8], Status: models.EventStatusActive,
	})
	if err != nil {
		t.Fatalf("seed event: %v", err)
	}
	fx.eventID = event.ID

	edition, err := r.Editions.Create(ctx, &models.Edition{
		EventID: event.ID, Name: "Activity Test Edition",
		Slug:     "activity-test-ed-" + uuid.NewString()[:8],
		StartsAt: time.Now().Add(-time.Hour), EndsAt: time.Now().Add(24 * time.Hour),
		CreatedBy: fx.ownerID,
	})
	if err != nil {
		t.Fatalf("seed edition: %v", err)
	}
	fx.editionID = edition.ID

	ticket, err := r.TicketTypes.Create(ctx, &models.TicketType{
		EditionID: edition.ID, Name: "Standard", AccessLevel: 0, PriceCents: 0,
	})
	if err != nil {
		t.Fatalf("seed ticket type: %v", err)
	}
	fx.ticketID = ticket.ID

	// Attendees A, B, C each hold a confirmed ticket.
	regs := map[uuid.UUID]*uuid.UUID{ // attendee -> reg id
		fx.attendeeA: &fx.regA,
		fx.attendeeB: &fx.regB,
		fx.attendeeC: &fx.regC,
	}
	for attendee, regID := range regs {
		reg, err := r.Registrations.Create(ctx, &models.Registration{
			EditionID: edition.ID, TicketTypeID: ticket.ID,
			PurchaserID: attendee, AttendeeUserID: &attendee,
			AttendeeEmail: "attendee-" + attendee.String() + "@example.com", AttendeeName: "Attendee",
			Status: models.RegistrationStatusConfirmed,
		})
		if err != nil {
			t.Fatalf("seed registration: %v", err)
		}
		*regID = reg.ID
	}

	program, err := r.Programs.Create(ctx, &models.Program{
		EditionID: edition.ID, Kind: models.ProgramKindActivity, Name: "Workshop",
		Price: new(int64(0)),
	})
	if err != nil {
		t.Fatalf("seed program: %v", err)
	}
	fx.programID = program.ID

	occurrence, err := r.Programs.CreateOccurrence(ctx, &models.ProgramOccurrence{
		ProgramID: program.ID, EditionID: edition.ID,
		StartsAt: time.Now().Add(2 * time.Hour), EndsAt: time.Now().Add(3 * time.Hour),
		MaxCapacity: new(int(2)),
	})
	if err != nil {
		t.Fatalf("seed occurrence: %v", err)
	}
	fx.occurrenceID = occurrence.ID

	_, err = r.Events.AddEventMember(ctx, event.ID, fx.staffID, models.EventMemberRoleStaff)
	if err != nil {
		t.Fatalf("seed staff member: %v", err)
	}

	return fx
}

func newOps(t *testing.T, r *repos.Repos, pool *pgxpool.Pool) (*programs.Operations, *recordingNotifier) {
	t.Helper()
	notifier := &recordingNotifier{}
	ops := programs.NewOperations(
		r.Events, r.Editions, r.Programs, r.Occurrences,
		r.Registrations, r.TicketTypes, r.Programs,
		authz.New(r.Events), notifier, database.NewPGXTxRunner(pool),
	)
	return ops, notifier
}

// newDB wires the real repos on a disposable Postgres running the real
// migrations, and the operations with the real tx runner.
func newDB(t *testing.T) (*repos.Repos, *programs.Operations, *recordingNotifier, *pgxpool.Pool) {
	t.Helper()
	pool := testdb.Postgres(t, "../../../db/migrations")
	r := repos.New(sqlc.New(pool))
	ops, notifier := newOps(t, r, pool)
	return r, ops, notifier, pool
}

// TestActivityRegistration_Lifecycle pins the full self-service loop on the
// real ledger: register → double-register 409 → capacity full 409 →
// deregister frees the slot → re-register appends a NEW row (append-only
// ledger) → staff marks attendance → the attendee's de-registration is
// locked (409).
func TestActivityRegistration_Lifecycle(t *testing.T) {
	r, ops, notifier, pool := newDB(t)
	fx := seedActivity(t, r)

	partA := mustRegister(t, ops, fx.occurrenceID, fx.attendeeA)
	if partA.Status != models.ProgramParticipationStatusRegistered {
		t.Fatalf("status = %s, want registered", partA.Status)
	}

	// A is already in → 409. B fills the second slot; C is then full → 409.
	wantRegisterConflict(t, ops, fx.occurrenceID, fx.attendeeA)
	mustRegister(t, ops, fx.occurrenceID, fx.attendeeB)
	wantRegisterConflict(t, ops, fx.occurrenceID, fx.attendeeC)

	// A deregisters → slot freed, then re-registers: the ledger keeps both
	// rows (append-only), the active one re-occupies the slot. C stays full.
	mustDeregister(t, ops, fx.occurrenceID, fx.attendeeA)
	mustRegister(t, ops, fx.occurrenceID, fx.attendeeA)
	assertLedgerRows(t, pool, fx.occurrenceID, fx.regA, 2)
	wantRegisterConflict(t, ops, fx.occurrenceID, fx.attendeeC)

	// Staff marks B attended → B's de-registration is locked (409); A
	// (still registered) can drop out.
	partB := participantFor(t, ops, fx.staffID, fx.occurrenceID, fx.regB)
	marked := mustMarkAttended(t, ops, partB.ID, fx.staffID)
	if marked.Status != models.ProgramParticipationStatusAttended {
		t.Fatalf("status = %s, want attended", marked.Status)
	}
	wantDeregisterConflict(t, ops, fx.occurrenceID, fx.attendeeB)
	mustDeregister(t, ops, fx.occurrenceID, fx.attendeeA)

	// Every write published a stock delta for the occurrence.
	for _, p := range notifier.stockPayloads() {
		if p.ItemIDs[0] != fx.occurrenceID {
			t.Fatalf("stock delta = %+v, want occurrence %s", p, fx.occurrenceID)
		}
	}
	if len(notifier.stockPayloads()) < 4 {
		t.Fatalf("want ≥4 stock deltas (register×2, deregister×2), got %d", len(notifier.stockPayloads()))
	}
}

// seedCheckpoint creates the checkpoint graph: event → edition → ticket
// type → confirmed registrations → checkpoint program (no price, no
// capacity — checkpoints are pass-through) → occurrence.
func seedCheckpoint(t *testing.T, r *repos.Repos) fixture {
	t.Helper()
	ctx := context.Background()
	fx := fixture{
		ownerID:   uuid.New(),
		attendeeA: uuid.New(),
		attendeeB: uuid.New(),
		staffID:   uuid.New(),
	}

	event, err := r.Events.Create(ctx, &models.Event{
		OwnerID: fx.ownerID, FullName: "Checkpoint Test Event",
		Slug: "checkpoint-test-" + uuid.NewString()[:8], Status: models.EventStatusActive,
	})
	if err != nil {
		t.Fatalf("seed event: %v", err)
	}
	fx.eventID = event.ID

	edition, err := r.Editions.Create(ctx, &models.Edition{
		EventID: event.ID, Name: "Checkpoint Test Edition",
		Slug:     "checkpoint-test-ed-" + uuid.NewString()[:8],
		StartsAt: time.Now().Add(-time.Hour), EndsAt: time.Now().Add(24 * time.Hour),
		CreatedBy: fx.ownerID,
	})
	if err != nil {
		t.Fatalf("seed edition: %v", err)
	}
	fx.editionID = edition.ID

	ticket, err := r.TicketTypes.Create(ctx, &models.TicketType{
		EditionID: edition.ID, Name: "Standard", AccessLevel: 0, PriceCents: 0,
	})
	if err != nil {
		t.Fatalf("seed ticket type: %v", err)
	}
	fx.ticketID = ticket.ID

	// Attendees A and B each hold a confirmed ticket.
	regs := map[uuid.UUID]*uuid.UUID{
		fx.attendeeA: &fx.regA,
		fx.attendeeB: &fx.regB,
	}
	for attendee, regID := range regs {
		reg, err := r.Registrations.Create(ctx, &models.Registration{
			EditionID: edition.ID, TicketTypeID: ticket.ID,
			PurchaserID: attendee, AttendeeUserID: &attendee,
			AttendeeEmail: "attendee-" + attendee.String() + "@example.com", AttendeeName: "Attendee",
			Status: models.RegistrationStatusConfirmed,
		})
		if err != nil {
			t.Fatalf("seed registration: %v", err)
		}
		*regID = reg.ID
	}

	program, err := r.Programs.Create(ctx, &models.Program{
		EditionID: edition.ID, Kind: models.ProgramKindCheckpoint, Name: "Security Gate",
	})
	if err != nil {
		t.Fatalf("seed checkpoint program: %v", err)
	}
	fx.checkpointProgramID = program.ID

	occurrence, err := r.Programs.CreateOccurrence(ctx, &models.ProgramOccurrence{
		ProgramID: program.ID, EditionID: edition.ID,
		StartsAt: time.Now().Add(2 * time.Hour), EndsAt: time.Now().Add(3 * time.Hour),
	})
	if err != nil {
		t.Fatalf("seed checkpoint occurrence: %v", err)
	}
	fx.checkpointOccurrenceID = occurrence.ID

	_, err = r.Events.AddEventMember(ctx, event.ID, fx.staffID, models.EventMemberRoleStaff)
	if err != nil {
		t.Fatalf("seed staff member: %v", err)
	}

	return fx
}

// TestCheckpointCheckIn_Lifecycle pins the checkpoint flow on the real
// ledger: the first scan upserts an attended participation anchored to the
// attendee's edition-ticket registration, a re-scan is idempotent (never
// duplicates — the partial unique index is the conflict target), the
// attendee shows up on the staff surface, and no stock delta is published
// (check-in does not move availability).
func TestCheckpointCheckIn_Lifecycle(t *testing.T) {
	r, ops, notifier, pool := newDB(t)
	fx := seedCheckpoint(t, r)

	part, err := ops.CheckIn(context.Background(), fx.checkpointOccurrenceID, fx.attendeeA, fx.staffID)
	if err != nil {
		t.Fatalf("check-in: %v", err)
	}
	if part.Status != models.ProgramParticipationStatusAttended {
		t.Fatalf("status = %s, want attended", part.Status)
	}
	if part.RegistrationID != fx.regA {
		t.Fatalf("registration = %s, want attendee A's ticket %s", part.RegistrationID, fx.regA)
	}
	assertLedgerRows(t, pool, fx.checkpointOccurrenceID, fx.regA, 1)

	// Re-scan: same attendee, still one row, still attended.
	again, err := ops.CheckIn(context.Background(), fx.checkpointOccurrenceID, fx.attendeeA, fx.staffID)
	if err != nil {
		t.Fatalf("re-check-in: %v", err)
	}
	if again.Status != models.ProgramParticipationStatusAttended {
		t.Fatalf("re-check-in status = %s, want attended", again.Status)
	}
	assertLedgerRows(t, pool, fx.checkpointOccurrenceID, fx.regA, 1)

	// The staff attendance surface now shows the checked-in attendee.
	list, err := ops.Participants(context.Background(), fx.checkpointOccurrenceID, fx.staffID)
	if err != nil {
		t.Fatalf("participants: %v", err)
	}
	if len(list) != 1 || list[0].RegistrationID != fx.regA {
		t.Fatalf("participants = %+v, want only attendee A", list)
	}

	// No stock delta: check-in never moves availability.
	if len(notifier.stockPayloads()) != 0 {
		t.Fatalf("stock payloads = %+v, want none", notifier.stockPayloads())
	}
}

// TestCheckpointCheckIn_Gates pins the pass gate on the real ledger: only
// event staff may check in, only ticket-holding attendees can pass, and
// only checkpoint occurrences accept check-in.
func TestCheckpointCheckIn_Gates(t *testing.T) {
	r, ops, _, _ := newDB(t)
	fx := seedCheckpoint(t, r)

	// Non-staff actor → 403.
	_, err := ops.CheckIn(context.Background(), fx.checkpointOccurrenceID, fx.attendeeA, fx.attendeeB)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("check-in as non-staff: want 403, got %v", err)
	}

	// Attendee without a ticket → 403 (cannot pass).
	_, err = ops.CheckIn(context.Background(), fx.checkpointOccurrenceID, uuid.New(), fx.staffID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("check-in without ticket: want 403, got %v", err)
	}

	// Activity occurrence → 400 (activities keep mark-attended).
	actFx := seedActivity(t, r)
	_, err = ops.CheckIn(context.Background(), actFx.occurrenceID, fx.attendeeA, actFx.staffID)
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("check-in on activity: want 400, got %v", err)
	}
}

// ── Lifecycle helpers ──────────────────────────────────────────────────────

func mustRegister(t *testing.T, ops *programs.Operations, occurrenceID, attendee uuid.UUID) *models.ProgramParticipation {
	t.Helper()
	part, err := ops.Register(context.Background(), occurrenceID, attendee)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	return part
}

func wantRegisterConflict(t *testing.T, ops *programs.Operations, occurrenceID, attendee uuid.UUID) {
	t.Helper()
	_, err := ops.Register(context.Background(), occurrenceID, attendee)
	if err == nil || !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("register: want 409, got %v", err)
	}
}

func mustDeregister(t *testing.T, ops *programs.Operations, occurrenceID, attendee uuid.UUID) {
	t.Helper()
	_, err := ops.Deregister(context.Background(), occurrenceID, attendee)
	if err != nil {
		t.Fatalf("deregister: %v", err)
	}
}

func wantDeregisterConflict(t *testing.T, ops *programs.Operations, occurrenceID, attendee uuid.UUID) {
	t.Helper()
	_, err := ops.Deregister(context.Background(), occurrenceID, attendee)
	if err == nil || !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("deregister: want 409, got %v", err)
	}
}

func wantMarkAttendedConflict(t *testing.T, ops *programs.Operations, partID, actorID uuid.UUID) {
	t.Helper()
	_, err := ops.MarkAttended(context.Background(), partID, actorID)
	if err == nil || !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("mark attended: want 409, got %v", err)
	}
}

// assertLedgerRows pins the append-only ledger: N rows exist for the
// (occurrence, registration) pair — cancelled rows are kept, re-registration
// inserts a fresh one.
func assertLedgerRows(t *testing.T, pool *pgxpool.Pool, occurrenceID, regID uuid.UUID, want int) {
	t.Helper()
	var rowCount int
	err := pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM program_participations WHERE occurrence_id = $1 AND registration_id = $2`,
		occurrenceID, regID).Scan(&rowCount)
	if err != nil {
		t.Fatalf("count ledger rows: %v", err)
	}
	if rowCount != want {
		t.Fatalf("ledger rows for (occ, reg) = %d, want %d", rowCount, want)
	}
}

// participantFor finds the participation of a registration in the staff
// attendance list.
func participantFor(t *testing.T, ops *programs.Operations, actorID, occurrenceID, regID uuid.UUID) models.ProgramParticipant {
	t.Helper()
	list, err := ops.Participants(context.Background(), occurrenceID, actorID)
	if err != nil {
		t.Fatalf("participants: %v", err)
	}
	for _, p := range list {
		if p.RegistrationID == regID {
			return p
		}
	}
	t.Fatalf("participant list missing registration %s: %+v", regID, list)
	return models.ProgramParticipant{}
}

func mustMarkAttended(t *testing.T, ops *programs.Operations, partID, actorID uuid.UUID) *models.ProgramParticipation {
	t.Helper()
	marked, err := ops.MarkAttended(context.Background(), partID, actorID)
	if err != nil {
		t.Fatalf("mark attended: %v", err)
	}
	return marked
}

// TestActivityRegistration_Availability pins the one-pool capacity rule: the
// store's availability ledger counts self-service registrations.
func TestActivityRegistration_Availability(t *testing.T) {
	r, ops, _, _ := newDB(t)
	fx := seedActivity(t, r)

	avail := func() int64 {
		t.Helper()
		items, err := r.Purchases.Availability(context.Background(), fx.editionID)
		if err != nil {
			t.Fatalf("availability: %v", err)
		}
		for _, it := range items {
			if it.ItemType == models.PurchaseItemTypeProgramOccurrence && it.ItemID == fx.occurrenceID {
				return it.ReservedQuantity
			}
		}
		t.Fatalf("occurrence %s missing from availability", fx.occurrenceID)
		return 0
	}

	if got := avail(); got != 0 {
		t.Fatalf("initial reserved = %d, want 0", got)
	}
	mustRegister(t, ops, fx.occurrenceID, fx.attendeeA)
	if got := avail(); got != 1 {
		t.Fatalf("reserved after register = %d, want 1", got)
	}
	mustDeregister(t, ops, fx.occurrenceID, fx.attendeeA)
	if got := avail(); got != 0 {
		t.Fatalf("reserved after deregister = %d, want 0", got)
	}
}

// TestActivityRegistration_MarkAttendedAuthz pins the staff gate: only an
// event staff member can mark attendance.
func TestActivityRegistration_MarkAttendedAuthz(t *testing.T) {
	r, ops, _, _ := newDB(t)
	fx := seedActivity(t, r)

	part := mustRegister(t, ops, fx.occurrenceID, fx.attendeeA)

	// Non-member → 403.
	_, err := ops.MarkAttended(context.Background(), part.ID, fx.attendeeB)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("mark as non-staff: want 403, got %v", err)
	}

	// Staff → attended, and re-marking stays idempotent.
	again := mustMarkAttended(t, ops, part.ID, fx.staffID)
	if again.Status != models.ProgramParticipationStatusAttended {
		t.Fatalf("status = %s, want attended", again.Status)
	}
	again = mustMarkAttended(t, ops, part.ID, fx.staffID)
	if again.Status != models.ProgramParticipationStatusAttended {
		t.Fatalf("re-mark status = %s, want attended", again.Status)
	}
}

// TestActivityRegistration_CancelledNotMarked pins the guard: a cancelled
// participation (dropped out / refunded) can never be marked attended.
func TestActivityRegistration_CancelledNotMarked(t *testing.T) {
	r, ops, _, _ := newDB(t)
	fx := seedActivity(t, r)

	part := mustRegister(t, ops, fx.occurrenceID, fx.attendeeA)
	mustDeregister(t, ops, fx.occurrenceID, fx.attendeeA)

	wantMarkAttendedConflict(t, ops, part.ID, fx.staffID)
}

// TestActivityRegistration_CertsExcludeCancelled pins the behavior change:
// a dropped-out attendee must not earn the program certificate.
func TestActivityRegistration_CertsExcludeCancelled(t *testing.T) {
	r, ops, _, _ := newDB(t)
	fx := seedActivity(t, r)

	// A registers then drops out; B stays registered.
	mustRegister(t, ops, fx.occurrenceID, fx.attendeeA)
	mustDeregister(t, ops, fx.occurrenceID, fx.attendeeA)
	mustRegister(t, ops, fx.occurrenceID, fx.attendeeB)

	participants, err := r.Certs.ListDistinctParticipantsByProgram(context.Background(), fx.programID)
	if err != nil {
		t.Fatalf("participants by program: %v", err)
	}
	if len(participants) != 1 || participants[0].UserID != fx.attendeeB {
		t.Fatalf("cert-eligible = %+v, want only attendee B", participants)
	}
}
