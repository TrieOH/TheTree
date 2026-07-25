package commands_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	idx "sdk/identityx"
	"univents/internal/features/events/commands"
	"univents/models"
	"univents/ports"
)

func TestDiscontinue_OwnerCanDiscontinueActive(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()

	eventID := uuid.New()
	ownerID := uuid.New()

	cmd := commands.NewCommands(repo, nil, nil)

	ctx := idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: ownerID},
	})

	event := &models.Event{
		ID:      eventID,
		OwnerID: ownerID,
		Status:  models.EventStatusActive,
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)
	mock.When(repo.Discontinue(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)

	err := cmd.Discontinue(ctx, eventID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_ = mock.Verify(repo, mock.Once()).Discontinue(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestDiscontinue_AdminCanDiscontinueActive(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()

	eventID := uuid.New()
	ownerID := uuid.New()
	adminID := uuid.New()

	cmd := commands.NewCommands(repo, nil, nil)

	ctx := idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: adminID},
	})

	event := &models.Event{
		ID:      eventID,
		OwnerID: ownerID,
		Status:  models.EventStatusActive,
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)
	mock.When(repo.GetMember(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(&models.EventMember{Role: models.EventMemberRoleAdmin}, nil)
	mock.When(repo.Discontinue(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)

	err := cmd.Discontinue(ctx, eventID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_ = mock.Verify(repo, mock.Once()).Discontinue(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestDiscontinue_CannotDiscontinueNonActive(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()

	eventID := uuid.New()
	ownerID := uuid.New()

	cmd := commands.NewCommands(repo, nil, nil)

	ctx := idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: ownerID},
	})

	event := &models.Event{
		ID:      eventID,
		OwnerID: ownerID,
		Status:  models.EventStatusDraft,
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)

	err := cmd.Discontinue(ctx, eventID)
	if err == nil {
		t.Fatal("expected bad request error, got nil")
	}

	_ = mock.Verify(repo, mock.Never()).Discontinue(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestDiscontinue_StaffForbidden(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()

	eventID := uuid.New()
	ownerID := uuid.New()
	staffID := uuid.New()

	cmd := commands.NewCommands(repo, nil, nil)

	ctx := idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: staffID},
	})

	event := &models.Event{
		ID:      eventID,
		OwnerID: ownerID,
		Status:  models.EventStatusActive,
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)
	mock.When(repo.GetMember(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(&models.EventMember{Role: models.EventMemberRoleStaff}, nil)

	err := cmd.Discontinue(ctx, eventID)
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}

	_ = mock.Verify(repo, mock.Never()).Discontinue(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestDiscontinue_NoIdentity(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()

	cmd := commands.NewCommands(repo, nil, nil)

	err := cmd.Discontinue(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected identity error, got nil")
	}
}
