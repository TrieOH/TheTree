package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"lib/database"
	"univents/internal/features/certifications/jobs"
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

	enq := &recordingJobs{}
	ops := newOps(programs, editions, authzEvents, enq)

	err := ops.EmitCertsForProgram(ownerCtx(), programID)
	if err != nil {
		t.Fatalf("EmitCertsForProgram: %v", err)
	}

	if len(enq.inserted) != 1 {
		t.Fatalf("want 1 inserted job, got %d", len(enq.inserted))
	}
	args, ok := enq.inserted[0].(jobs.GrantCertsForOccurrenceArgs)
	if !ok {
		t.Fatalf("want GrantCertsForOccurrenceArgs, got %T", enq.inserted[0])
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

	enq := &recordingJobs{}
	ops := newOps(programs, editions, authzEvents, enq)

	err := ops.EmitCertsForProgram(ownerCtx(), programID)
	if err == nil {
		t.Fatal("want forbidden error, got nil")
	}
	if len(enq.inserted) != 0 {
		t.Fatal("job must not be enqueued for a non-admin actor")
	}
}

func TestEmitCertsForProgram_UnknownProgram(t *testing.T) {
	mock.SetUp(t)

	programs := mock.Mock[ports.ProgramRepo]()
	mock.When(programs.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(nil, database.ErrNotFound)

	enq := &recordingJobs{}
	ops := newOps(programs, mock.Mock[ports.EditionRepo](), mock.Mock[ports.EventRepo](), enq)

	err := ops.EmitCertsForProgram(ownerCtx(), uuid.New())
	if err == nil {
		t.Fatal("want not-found error, got nil")
	}
	if len(enq.inserted) != 0 {
		t.Fatal("job must not be enqueued for an unknown program")
	}
}
