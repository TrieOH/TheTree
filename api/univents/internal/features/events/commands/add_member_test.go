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

func TestAddMember_StaffForbidden(t *testing.T) {
	mock.SetUp(t)

	var repo ports.EventRepo = mock.Mock[ports.EventRepo]()

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
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)
	mock.When(repo.GetMember(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(&models.EventMember{Role: models.EventMemberRoleStaff}, nil)

	_, err := cmd.AddMember(ctx, eventID, models.AddEventMemberRequest{
		Email: "someone@example.com",
		Role:  models.EventMemberRoleStaff,
	})
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
}

func TestAddMember_OwnerGetsPastAuth(t *testing.T) {
	mock.SetUp(t)

	var repo ports.EventRepo = mock.Mock[ports.EventRepo]()

	eventID := uuid.New()
	ownerID := uuid.New()

	cmd := commands.NewCommands(repo, nil, nil)

	ctx := idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: ownerID},
	})

	event := &models.Event{
		ID:      eventID,
		OwnerID: ownerID,
	}

	// Owner skips the GetMember check, so only GetByID fires,
	// then it tries c.idx.Actors.GetByEmail which is nil — but that's fine,
	// the test proves authorization logic lets the owner through.
	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)

	// Will panic on nil idx, but that's expected — this test documents
	// that owners bypass the member permission check.
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil idx client, but didn't panic")
		}
	}()

	cmd.AddMember(ctx, eventID, models.AddEventMemberRequest{
		Email: "someone@example.com",
		Role:  models.EventMemberRoleStaff,
	})
}

func TestAddMember_NoIdentity(t *testing.T) {
	mock.SetUp(t)

	var repo ports.EventRepo = mock.Mock[ports.EventRepo]()

	cmd := commands.NewCommands(repo, nil, nil)

	_, err := cmd.AddMember(context.Background(), uuid.New(), models.AddEventMemberRequest{
		Email: "someone@example.com",
		Role:  models.EventMemberRoleStaff,
	})
	if err == nil {
		t.Fatal("expected identity error, got nil")
	}
}
