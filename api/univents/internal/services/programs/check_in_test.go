package programs_test

import (
	"context"
	"testing"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"univents/models"
)

func checkpointProgram(programID uuid.UUID, minAccess *int, staffOnly bool) *models.Program {
	return &models.Program{
		ID: programID, Kind: models.ProgramKindCheckpoint, Name: "Security Gate",
		MinAccessLevel: minAccess, StaffOnly: staffOnly,
	}
}

// stubCheckInStaff wires the staff-authz mocks for the checkpoint flow:
// the occurrence resolves, its edition resolves, and the actor holds an
// event staff role. Returns the edition so tests can match the attendee's
// role check.
func (h *harness) stubCheckInStaff(t *testing.T, occ *models.ProgramOccurrence, staffID uuid.UUID) *models.Edition {
	t.Helper()
	edition := &models.Edition{ID: occ.EditionID, EventID: uuid.New()}
	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(occ.EditionID))).ThenReturn(edition, nil)
	mock.When(h.events.GetRole(mock.AnyContext(), mock.Equal(staffID), mock.Equal(edition.EventID))).ThenReturn(models.EventMemberRoleStaff, nil)
	return edition
}

func TestCheckIn_HappyPath(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	prog := checkpointProgram(occ.ProgramID, nil, false)
	staffID := uuid.New()
	h.stubCheckInStaff(t, occ, staffID)

	mock.When(h.programs.GetByID(mock.AnyContext(), mock.Equal(occ.ProgramID))).ThenReturn(prog, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(*reg.AttendeeUserID))).ThenReturn(reg, nil)
	checkedIn := &models.ProgramParticipation{ID: uuid.New(), EditionID: occ.EditionID, OccurrenceID: occ.ID, RegistrationID: reg.ID, Status: models.ProgramParticipationStatusAttended}
	mock.When(h.participations.UpsertAttended(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(occ.ID), mock.Equal(reg.ID))).ThenReturn(checkedIn, nil)

	got, err := h.ops.CheckIn(context.Background(), occ.ID, *reg.AttendeeUserID, staffID)
	if err != nil {
		t.Fatalf("CheckIn: %v", err)
	}
	if got.Status != models.ProgramParticipationStatusAttended {
		t.Fatalf("status = %s, want attended", got.Status)
	}
	if got.RegistrationID != reg.ID {
		t.Fatalf("registration = %s, want the attendee's edition ticket %s", got.RegistrationID, reg.ID)
	}
	// Check-in never publishes a stock delta (availability unchanged).
	if len(h.notifier.stockPayloads()) != 0 {
		t.Fatalf("stock payloads = %+v, want none", h.notifier.stockPayloads())
	}
}

func TestCheckIn_NonCheckpoint400(t *testing.T) {
	// Activities keep the participation-driven mark-attended flow.
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	prog := program(occ.ProgramID, nil, false) // activity
	staffID := uuid.New()
	h.stubCheckInStaff(t, occ, staffID)

	mock.When(h.programs.GetByID(mock.AnyContext(), mock.Equal(occ.ProgramID))).ThenReturn(prog, nil)

	_, err := h.ops.CheckIn(context.Background(), occ.ID, uuid.New(), staffID)
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("want 400, got %v", err)
	}
}

func TestCheckIn_NonStaff403(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	prog := checkpointProgram(occ.ProgramID, nil, false)
	userID := uuid.New()
	edition := &models.Edition{ID: occ.EditionID, EventID: uuid.New()}

	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occ.ID))).ThenReturn(occ, nil)
	mock.When(h.editions.GetByID(mock.AnyContext(), mock.Equal(occ.EditionID))).ThenReturn(edition, nil)
	mock.When(h.events.GetRole(mock.AnyContext(), mock.Equal(userID), mock.Equal(edition.EventID))).ThenReturn(models.EventMemberRole(""), notFound())
	mock.When(h.programs.GetByID(mock.AnyContext(), mock.Equal(occ.ProgramID))).ThenReturn(prog, nil)

	_, err := h.ops.CheckIn(context.Background(), occ.ID, uuid.New(), userID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}

func TestCheckIn_NoTicket403(t *testing.T) {
	// The scanned attendee holds no edition ticket: cannot pass.
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	prog := checkpointProgram(occ.ProgramID, nil, false)
	staffID := uuid.New()
	h.stubCheckInStaff(t, occ, staffID)
	attendeeID := uuid.New()

	mock.When(h.programs.GetByID(mock.AnyContext(), mock.Equal(occ.ProgramID))).ThenReturn(prog, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(attendeeID))).ThenReturn(nil, notFound())

	_, err := h.ops.CheckIn(context.Background(), occ.ID, attendeeID, staffID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}

func TestCheckIn_PendingTicket403(t *testing.T) {
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	reg.Status = models.RegistrationStatusPending
	prog := checkpointProgram(occ.ProgramID, nil, false)
	staffID := uuid.New()
	h.stubCheckInStaff(t, occ, staffID)

	mock.When(h.programs.GetByID(mock.AnyContext(), mock.Equal(occ.ProgramID))).ThenReturn(prog, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(*reg.AttendeeUserID))).ThenReturn(reg, nil)

	_, err := h.ops.CheckIn(context.Background(), occ.ID, *reg.AttendeeUserID, staffID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}

func TestCheckIn_LowAccessLevel403(t *testing.T) {
	// A checkpoint with min_access_level gates on the attendee's ticket.
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	prog := checkpointProgram(occ.ProgramID, new(int(2)), false)
	staffID := uuid.New()
	h.stubCheckInStaff(t, occ, staffID)

	mock.When(h.programs.GetByID(mock.AnyContext(), mock.Equal(occ.ProgramID))).ThenReturn(prog, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(*reg.AttendeeUserID))).ThenReturn(reg, nil)
	mock.When(h.ticketTypes.GetByID(mock.AnyContext(), mock.Equal(reg.TicketTypeID))).ThenReturn(&models.TicketType{ID: reg.TicketTypeID, AccessLevel: 1}, nil)

	_, err := h.ops.CheckIn(context.Background(), occ.ID, *reg.AttendeeUserID, staffID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}

func TestCheckIn_StaffOnlyNonStaff403(t *testing.T) {
	// A staff-only checkpoint: the attendee must hold an event staff role.
	h := newHarness(t)
	occ := occurrence(uuid.New(), uuid.New(), nil)
	reg := confirmedRegistration(uuid.New(), occ.EditionID)
	prog := checkpointProgram(occ.ProgramID, nil, true)
	staffID := uuid.New()
	edition := h.stubCheckInStaff(t, occ, staffID)

	mock.When(h.programs.GetByID(mock.AnyContext(), mock.Equal(occ.ProgramID))).ThenReturn(prog, nil)
	mock.When(h.registrations.GetActiveByEditionAndAttendee(mock.AnyContext(), mock.Equal(occ.EditionID), mock.Equal(*reg.AttendeeUserID))).ThenReturn(reg, nil)
	// Attendee is not staff.
	mock.When(h.events.GetRole(mock.AnyContext(), mock.Equal(*reg.AttendeeUserID), mock.Equal(edition.EventID))).ThenReturn(models.EventMemberRole(""), notFound())

	_, err := h.ops.CheckIn(context.Background(), occ.ID, *reg.AttendeeUserID, staffID)
	if err == nil || !fun.Is(err, fun.CodeForbidden) {
		t.Fatalf("want 403, got %v", err)
	}
}

func TestCheckIn_UnknownOccurrence404(t *testing.T) {
	h := newHarness(t)
	occID := uuid.New()

	mock.When(h.occurrences.GetOccurrenceByID(mock.AnyContext(), mock.Equal(occID))).ThenReturn(nil, notFound())

	_, err := h.ops.CheckIn(context.Background(), occID, uuid.New(), uuid.New())
	if err == nil || !fun.Is(err, fun.CodeNotFound) {
		t.Fatalf("want 404, got %v", err)
	}
}
