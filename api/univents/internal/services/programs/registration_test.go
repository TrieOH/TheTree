package programs_test

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"lib/database"

	"univents/internal/authz"
	"univents/internal/services/programs"
	"univents/models"
	"univents/ports"
)

func TestMain(m *testing.M) {
	fun.SetConfig(fun.Config{
		DefaultModule:        "test",
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
	})
	m.Run()
}

// noopTxRunner runs fn directly — the repo fakes ignore the tx context, so
// the tx is simulated. The DB-backed test exercises the real pgx runner.
type noopTxRunner struct{}

func (noopTxRunner) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
func (noopTxRunner) WithinTxWithOptions(ctx context.Context, _ database.TxOptions, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type notifyCall struct {
	channel string
	payload string
}

// recordingNotifier captures Notify calls for assertions.
type recordingNotifier struct {
	mu    sync.Mutex
	calls []notifyCall
}

func (n *recordingNotifier) Notify(_ context.Context, channel, payload string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.calls = append(n.calls, notifyCall{channel, payload})
	return nil
}

type decodedNotify struct {
	Kind      string      `json:"kind"`
	EditionID uuid.UUID   `json:"edition_id"`
	ItemIDs   []uuid.UUID `json:"item_ids"`
}

func (n *recordingNotifier) stockPayloads() []decodedNotify {
	n.mu.Lock()
	defer n.mu.Unlock()
	out := make([]decodedNotify, 0, len(n.calls))
	for _, c := range n.calls {
		var d decodedNotify
		err := json.Unmarshal([]byte(c.payload), &d)
		if err != nil {
			panic("bad notify payload: " + err.Error())
		}
		out = append(out, d)
	}
	return out
}

func notFound() error {
	return fun.Err("not found").NotFound()
}

type harness struct {
	events         ports.EventRepo
	editions       ports.EditionRepo
	programs       ports.ProgramRepo
	occurrences    ports.ProgramOccurrenceRepo
	registrations  ports.RegistrationRepo
	ticketTypes    ports.TicketTypeRepo
	participations ports.ProgramParticipationRepo
	notifier       *recordingNotifier
	ops            *programs.Operations
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	mock.SetUp(t)
	h := &harness{
		events:         mock.Mock[ports.EventRepo](),
		editions:       mock.Mock[ports.EditionRepo](),
		programs:       mock.Mock[ports.ProgramRepo](),
		occurrences:    mock.Mock[ports.ProgramOccurrenceRepo](),
		registrations:  mock.Mock[ports.RegistrationRepo](),
		ticketTypes:    mock.Mock[ports.TicketTypeRepo](),
		participations: mock.Mock[ports.ProgramParticipationRepo](),
		notifier:       &recordingNotifier{},
	}
	h.ops = programs.NewOperations(
		h.events, h.editions, h.programs, h.occurrences,
		h.registrations, h.ticketTypes, h.participations,
		authz.New(h.events), h.notifier, noopTxRunner{},
	)
	return h
}

func occurrence(id uuid.UUID, programID uuid.UUID, capacity *int) *models.ProgramOccurrence {
	return &models.ProgramOccurrence{
		ID: id, ProgramID: programID, EditionID: uuid.New(),
		StartsAt: time.Now().Add(time.Hour), EndsAt: time.Now().Add(2 * time.Hour), MaxCapacity: capacity,
	}
}

func program(programID uuid.UUID, minAccess *int, staffOnly bool) *models.Program {
	return &models.Program{
		ID: programID, Kind: models.ProgramKindActivity, Name: "Workshop",
		MinAccessLevel: minAccess, StaffOnly: staffOnly,
	}
}

func confirmedRegistration(regID, editionID uuid.UUID) *models.Registration {
	return &models.Registration{
		ID: regID, EditionID: editionID, TicketTypeID: uuid.New(),
		AttendeeUserID: new(uuid.New()), Status: models.RegistrationStatusConfirmed,
	}
}

// stubEligible wires the common happy-path mocks: occurrence lookup, the
// caller's confirmed registration, and the program (open to all).
func (h *harness) stubEligible(t *testing.T, occ *models.ProgramOccurrence, reg *models.Registration, prog *models.Program) {
	t.Helper()
	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(*reg.AttendeeUserID))).ThenReturn(reg, nil)
	mock.When(h.programs.GetByID(mock.AnyContext(), mock.Equal(occ.ProgramID))).ThenReturn(prog, nil)
}

func TestRegister_HappyPath(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), new(int(10)))
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	prog := program(occ.ProgramID, nil, false)
	h.stubEligible(t, occ, reg, prog)

	mock.When(h.occurrences.GetOccurrenceByIDForUpdate(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.participations.GetActiveByOccurrenceAndRegistration(mock.AnyContext(), mock.Equal(occ.ID), mock.Equal(reg.ID))).ThenReturn(nil, notFound())
	mock.When(h.participations.CountActiveByOccurrence(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(int64(3), nil)
	created := &models.ProgramParticipation{ID: uuid.New(), EditionID: occ.EditionID, OccurrenceID: occ.ID, RegistrationID: reg.ID, Status: models.ProgramParticipationStatusRegistered}
	mock.When(h.participations.CreateParticipation(mock.AnyContext(), mock.Any[*models.ProgramParticipation]())).ThenReturn(created, nil)

	part, err := h.ops.Register(context.Background(), occ.ID, *reg.AttendeeUserID)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if part.ID != created.ID {
		t.Fatalf("part = %+v, want created row", part)
	}
	// Stock delta published for the occurrence.
	payloads := h.notifier.stockPayloads()
	if len(payloads) != 1 || payloads[0].ItemIDs[0] != occ.ID {
		t.Fatalf("stock payloads = %+v, want one delta for %s", payloads, occ.ID)
	}
}

func TestRegister_NoTicket403(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), new(int(10)))
	userID := uuid.New()

	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(userID))).ThenReturn(nil, notFound())

	_, err := h.ops.Register(context.Background(), occ.ID, userID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}

func TestRegister_PendingTicket403(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), new(int(10)))
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	reg.Status = models.RegistrationStatusPending

	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(*reg.AttendeeUserID))).ThenReturn(reg, nil)

	_, err := h.ops.Register(context.Background(), occ.ID, *reg.AttendeeUserID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}

func TestRegister_LowAccessLevel403(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), new(int(10)))
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	prog := program(occ.ProgramID, new(int(2)), false)
	h.stubEligible(t, occ, reg, prog)

	mock.When(h.ticketTypes.GetByID(mock.AnyContext(), mock.Equal(reg.TicketTypeID))).ThenReturn(&models.TicketType{ID: reg.TicketTypeID, AccessLevel: 1}, nil)

	_, err := h.ops.Register(context.Background(), occ.ID, *reg.AttendeeUserID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}

func TestRegister_StaffOnlyNonStaff403(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), new(int(10)))
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	prog := program(occ.ProgramID, nil, true)
	h.stubEligible(t, occ, reg, prog)

	edition := &models.Edition{ID: occ.EditionID, EventID: uuid.New()}
	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(occ.EditionID))).ThenReturn(edition, nil)
	// Not a member → authz maps NOT_FOUND to Forbidden.
	mock.When(h.events.GetRole(mock.AnyContext(), mock.Equal(*reg.AttendeeUserID), mock.Equal(edition.EventID))).ThenReturn(models.EventMemberRole(""), notFound())

	_, err := h.ops.Register(context.Background(), occ.ID, *reg.AttendeeUserID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}

func TestRegister_StaffOnlyStaffOK(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), new(int(10)))
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	prog := program(occ.ProgramID, nil, true)
	h.stubEligible(t, occ, reg, prog)

	edition := &models.Edition{ID: occ.EditionID, EventID: uuid.New()}
	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(occ.EditionID))).ThenReturn(edition, nil)
	mock.When(h.events.GetRole(mock.AnyContext(), mock.Equal(*reg.AttendeeUserID), mock.Equal(edition.EventID))).ThenReturn(models.EventMemberRoleStaff, nil)
	mock.When(h.occurrences.GetOccurrenceByIDForUpdate(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.participations.GetActiveByOccurrenceAndRegistration(mock.AnyContext(), mock.Equal(occ.ID), mock.Equal(reg.ID))).ThenReturn(nil, notFound())
	mock.When(h.participations.CountActiveByOccurrence(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(int64(0), nil)
	mock.When(h.participations.CreateParticipation(mock.AnyContext(), mock.Any[*models.ProgramParticipation]())).ThenReturn(&models.ProgramParticipation{ID: uuid.New(), Status: models.ProgramParticipationStatusRegistered}, nil)

	_, err := h.ops.Register(context.Background(), occ.ID, *reg.AttendeeUserID)
	if err != nil {
		t.Fatalf("Register (staff): %v", err)
	}
}

func TestRegister_AlreadyRegistered409(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), new(int(10)))
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	prog := program(occ.ProgramID, nil, false)
	h.stubEligible(t, occ, reg, prog)

	mock.When(h.occurrences.GetOccurrenceByIDForUpdate(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.participations.GetActiveByOccurrenceAndRegistration(mock.AnyContext(), mock.Equal(occ.ID), mock.Equal(reg.ID))).ThenReturn(&models.ProgramParticipation{ID: uuid.New()}, nil)

	_, err := h.ops.Register(context.Background(), occ.ID, *reg.AttendeeUserID)
	if err == nil || !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("want 409, got %v", err)
	}
}

func TestRegister_Full409(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), new(int(2)))
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	prog := program(occ.ProgramID, nil, false)
	h.stubEligible(t, occ, reg, prog)

	mock.When(h.occurrences.GetOccurrenceByIDForUpdate(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.participations.GetActiveByOccurrenceAndRegistration(mock.AnyContext(), mock.Equal(occ.ID), mock.Equal(reg.ID))).ThenReturn(nil, notFound())
	mock.When(h.participations.CountActiveByOccurrence(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(int64(2), nil)

	_, err := h.ops.Register(context.Background(), occ.ID, *reg.AttendeeUserID)
	if err == nil || !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("want 409, got %v", err)
	}
}

func TestRegister_UnknownOccurrence404(t *testing.T) {
	h := newHarness(t)
	occID := uuid.New()

	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occID))).ThenReturn(nil, notFound())

	_, err := h.ops.Register(context.Background(), occID, uuid.New())
	if err == nil || !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("want 404, got %v", err)
	}
}

func TestDeregister_HappyPath(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	part := &models.ProgramParticipation{ID: uuid.New(), EditionID: occ.EditionID, OccurrenceID: occ.ID, RegistrationID: reg.ID, Status: models.ProgramParticipationStatusRegistered}

	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(*reg.AttendeeUserID))).ThenReturn(reg, nil)
	mock.When(h.participations.GetActiveByOccurrenceAndRegistration(mock.AnyContext(), mock.Equal(occ.ID), mock.Equal(reg.ID))).ThenReturn(part, nil)
	cancelled := *part
	cancelled.Status = models.ProgramParticipationStatusCancelled
	mock.When(h.participations.UpdateParticipationStatusIfRegistered(mock.AnyContext(), mock.Equal(part.ID), mock.Equal(models.ProgramParticipationStatusCancelled))).ThenReturn(&cancelled, nil)

	got, err := h.ops.Deregister(context.Background(), occ.ID, *reg.AttendeeUserID)
	if err != nil {
		t.Fatalf("Deregister: %v", err)
	}
	if got.Status != models.ProgramParticipationStatusCancelled {
		t.Fatalf("status = %s, want cancelled", got.Status)
	}
	if len(h.notifier.stockPayloads()) != 1 {
		t.Fatalf("want one stock delta, got %d", len(h.notifier.stockPayloads()))
	}
}

func TestDeregister_NotRegistered404(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	reg := confirmedRegistration(uuid.New(), occ.EditionID)

	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(*reg.AttendeeUserID))).ThenReturn(reg, nil)
	mock.When(h.participations.GetActiveByOccurrenceAndRegistration(mock.AnyContext(), mock.Equal(occ.ID), mock.Equal(reg.ID))).ThenReturn(nil, notFound())

	_, err := h.ops.Deregister(context.Background(), occ.ID, *reg.AttendeeUserID)
	if err == nil || !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("want 404, got %v", err)
	}
}

func TestDeregister_LockedWhenAttended409(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	part := &models.ProgramParticipation{ID: uuid.New(), EditionID: occ.EditionID, OccurrenceID: occ.ID, RegistrationID: reg.ID, Status: models.ProgramParticipationStatusAttended}

	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(*reg.AttendeeUserID))).ThenReturn(reg, nil)
	mock.When(h.participations.GetActiveByOccurrenceAndRegistration(mock.AnyContext(), mock.Equal(occ.ID), mock.Equal(reg.ID))).ThenReturn(part, nil)
	// The guarded update returns 0 rows — the lock.
	mock.When(h.participations.UpdateParticipationStatusIfRegistered(mock.AnyContext(), mock.Equal(part.ID), mock.Equal(models.ProgramParticipationStatusCancelled))).ThenReturn(nil, nil)

	_, err := h.ops.Deregister(context.Background(), occ.ID, *reg.AttendeeUserID)
	if err == nil || !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("want 409, got %v", err)
	}
}

func TestMarkAttended_HappyPath(t *testing.T) {
	h := newHarness(t)
	part := &models.ProgramParticipation{ID: uuid.New(), EditionID: uuid.New(), OccurrenceID: uuid.New(), RegistrationID: uuid.New(), Status: models.ProgramParticipationStatusRegistered}
	edition := &models.Edition{ID: part.EditionID, EventID: uuid.New()}
	staffID := uuid.New()

	mock.When(h.participations.GetParticipationByID(mock.AnyContext(), mock.Equal(part.ID))).ThenReturn(part, nil)
	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(part.EditionID))).ThenReturn(edition, nil)
	mock.When(h.events.GetRole(mock.AnyContext(), mock.Equal(staffID), mock.Equal(edition.EventID))).ThenReturn(models.EventMemberRoleStaff, nil)
	attended := *part
	attended.Status = models.ProgramParticipationStatusAttended
	mock.When(h.participations.UpdateParticipationStatusIfNotCancelled(mock.AnyContext(), mock.Equal(part.ID), mock.Equal(models.ProgramParticipationStatusAttended))).ThenReturn(&attended, nil)

	got, err := h.ops.MarkAttended(context.Background(), part.ID, staffID)
	if err != nil {
		t.Fatalf("MarkAttended: %v", err)
	}
	if got.Status != models.ProgramParticipationStatusAttended {
		t.Fatalf("status = %s, want attended", got.Status)
	}
}

func TestMarkAttended_NonStaff403(t *testing.T) {
	h := newHarness(t)
	part := &models.ProgramParticipation{ID: uuid.New(), EditionID: uuid.New(), OccurrenceID: uuid.New(), RegistrationID: uuid.New(), Status: models.ProgramParticipationStatusRegistered}
	edition := &models.Edition{ID: part.EditionID, EventID: uuid.New()}
	userID := uuid.New()

	mock.When(h.participations.GetParticipationByID(mock.AnyContext(), mock.Equal(part.ID))).ThenReturn(part, nil)
	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(part.EditionID))).ThenReturn(edition, nil)
	mock.When(h.events.GetRole(mock.AnyContext(), mock.Equal(userID), mock.Equal(edition.EventID))).ThenReturn(models.EventMemberRole(""), notFound())

	_, err := h.ops.MarkAttended(context.Background(), part.ID, userID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}

func TestMarkAttended_Cancelled409(t *testing.T) {
	h := newHarness(t)
	part := &models.ProgramParticipation{ID: uuid.New(), EditionID: uuid.New(), OccurrenceID: uuid.New(), RegistrationID: uuid.New(), Status: models.ProgramParticipationStatusCancelled}
	edition := &models.Edition{ID: part.EditionID, EventID: uuid.New()}
	staffID := uuid.New()

	mock.When(h.participations.GetParticipationByID(mock.AnyContext(), mock.Equal(part.ID))).ThenReturn(part, nil)
	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(part.EditionID))).ThenReturn(edition, nil)
	mock.When(h.events.GetRole(mock.AnyContext(), mock.Equal(staffID), mock.Equal(edition.EventID))).ThenReturn(models.EventMemberRoleStaff, nil)

	_, err := h.ops.MarkAttended(context.Background(), part.ID, staffID)
	if err == nil || !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("want 409, got %v", err)
	}
}

func TestMarkAttended_Unknown404(t *testing.T) {
	h := newHarness(t)
	partID := uuid.New()

	mock.When(h.participations.GetParticipationByID(mock.AnyContext(), mock.Equal(partID))).ThenReturn(nil, notFound())

	_, err := h.ops.MarkAttended(context.Background(), partID, uuid.New())
	if err == nil || !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("want 404, got %v", err)
	}
}

func TestMyParticipations_HappyPath(t *testing.T) {
	h := newHarness(t)
	editionID := uuid.New()
	reg := confirmedRegistration(uuid.New(), editionID)
	row := models.MyParticipation{ID: uuid.New(), EditionID: editionID, Status: models.ProgramParticipationStatusRegistered}

	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(editionID))).ThenReturn(&models.Edition{ID: editionID}, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(editionID), mock.Equal(*reg.AttendeeUserID))).ThenReturn(reg, nil)
	mock.When(h.participations.ListActiveByEditionAndRegistration(mock.AnyContext(), mock.Equal(editionID), mock.Equal(reg.ID))).ThenReturn([]models.MyParticipation{row}, nil)

	rows, err := h.ops.MyParticipations(context.Background(), editionID, *reg.AttendeeUserID)
	if err != nil {
		t.Fatalf("MyParticipations: %v", err)
	}
	if len(rows) != 1 || rows[0].ID != row.ID {
		t.Fatalf("rows = %+v, want the stubbed row", rows)
	}
}

func TestMyParticipations_NoTicketNil(t *testing.T) {
	h := newHarness(t)
	editionID := uuid.New()
	userID := uuid.New()

	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(editionID))).ThenReturn(&models.Edition{ID: editionID}, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(editionID), mock.Equal(userID))).ThenReturn(nil, notFound())

	rows, err := h.ops.MyParticipations(context.Background(), editionID, userID)
	if err != nil {
		t.Fatalf("MyParticipations: %v", err)
	}
	if rows != nil {
		t.Fatalf("rows = %+v, want nil", rows)
	}
}

func TestParticipants_HappyPath(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	edition := &models.Edition{ID: occ.EditionID, EventID: uuid.New()}
	staffID := uuid.New()
	row := models.ProgramParticipant{ID: uuid.New(), OccurrenceID: occ.ID, Status: models.ProgramParticipationStatusRegistered}

	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(occ.EditionID))).ThenReturn(edition, nil)
	mock.When(h.events.GetRole(mock.AnyContext(), mock.Equal(staffID), mock.Equal(edition.EventID))).ThenReturn(models.EventMemberRoleStaff, nil)
	mock.When(h.participations.ListByOccurrence(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn([]models.ProgramParticipant{row}, nil)

	rows, err := h.ops.Participants(context.Background(), occ.ID, staffID)
	if err != nil {
		t.Fatalf("Participants: %v", err)
	}
	if len(rows) != 1 || rows[0].ID != row.ID {
		t.Fatalf("rows = %+v, want the stubbed row", rows)
	}
}

func TestParticipants_NonStaff403(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	edition := &models.Edition{ID: occ.EditionID, EventID: uuid.New()}
	userID := uuid.New()

	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(occ.EditionID))).ThenReturn(edition, nil)
	mock.When(h.events.GetRole(mock.AnyContext(), mock.Equal(userID), mock.Equal(edition.EventID))).ThenReturn(models.EventMemberRole(""), notFound())

	_, err := h.ops.Participants(context.Background(), occ.ID, userID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}
