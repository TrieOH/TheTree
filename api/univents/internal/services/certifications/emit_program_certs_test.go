package certifications_test

import (
	"testing"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"univents/internal/services/certifications/jobs"
	"univents/models"
	"univents/ports"
)

func TestEmitCertsForProgram_EnqueuesGrantJob(t *testing.T) {
	mock.SetUp(t)

	programID := uuid.New()
	editionID := uuid.New()

	programs := mock.Mock[ports.ProgramRepo]()
	editions := mock.Mock[ports.EditionRepo]()
	authzEvents := mock.Mock[ports.EventRepo]()

	mock.When(programs.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Program{ID: programID, EditionID: editionID}, nil)
	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: editionID, EventID: uuid.New()}, nil)
	mock.When(authzEvents.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleAdmin, nil)

	river := &recordingRiver{}
	ops := newOps(programs, editions, authzEvents, river)

	err := ops.EmitCertsForProgram(ownerCtx(), programID)
	if err != nil {
		t.Fatalf("EmitCertsForProgram: %v", err)
	}

	if len(river.inserted) != 1 {
		t.Fatalf("want 1 inserted job, got %d", len(river.inserted))
	}
	args, ok := river.inserted[0].(jobs.GrantCertsForOccurrenceArgs)
	if !ok {
		t.Fatalf("want GrantCertsForOccurrenceArgs, got %T", river.inserted[0])
	}
	if args.EditionID != editionID || args.ProgramID != programID {
		t.Errorf("unexpected job args: %+v", args)
	}
}

func TestEmitCertsForProgram_StaffIsForbidden(t *testing.T) {
	mock.SetUp(t)

	programID := uuid.New()

	programs := mock.Mock[ports.ProgramRepo]()
	editions := mock.Mock[ports.EditionRepo]()
	authzEvents := mock.Mock[ports.EventRepo]()

	mock.When(programs.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Program{ID: programID, EditionID: uuid.New()}, nil)
	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: uuid.New(), EventID: uuid.New()}, nil)
	mock.When(authzEvents.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleStaff, nil)

	river := &recordingRiver{}
	ops := newOps(programs, editions, authzEvents, river)

	err := ops.EmitCertsForProgram(ownerCtx(), programID)
	if err == nil {
		t.Fatal("want forbidden error, got nil")
	}
	if len(river.inserted) != 0 {
		t.Fatal("job must not be enqueued for a non-admin actor")
	}
}

func TestEmitCertsForProgram_UnknownProgram(t *testing.T) {
	mock.SetUp(t)

	programs := mock.Mock[ports.ProgramRepo]()
	mock.When(programs.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(nil, fun.ErrNotFound("program not found"))

	river := &recordingRiver{}
	ops := newOps(programs, mock.Mock[ports.EditionRepo](), mock.Mock[ports.EventRepo](), river)

	err := ops.EmitCertsForProgram(ownerCtx(), uuid.New())
	if err == nil {
		t.Fatal("want not-found error, got nil")
	}
	if len(river.inserted) != 0 {
		t.Fatal("job must not be enqueued for an unknown program")
	}
}
