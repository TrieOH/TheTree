package events_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/internal/features/events"
	"univents/models"
	"univents/ports"
)

func TestRemoveMember_StaffForbidden(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()
	authz.Service = authz.New(repo)

	eventID := uuid.New()
	ownerID := uuid.New()
	staffID := uuid.New()

	cmd := events.NewOperations(repo, nil, nil)

	ctx := idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: staffID},
	})

	event := &models.Event{
		ID:      eventID,
		OwnerID: ownerID,
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)
	mock.When(repo.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleStaff, nil)

	err := cmd.RemoveMember(ctx, eventID, models.RemoveMemberRequest{
		Email: "someone@example.com",
	})
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
}

func TestRemoveMember_OwnerGetsPastAuth(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()
	authz.Service = authz.New(repo)

	eventID := uuid.New()
	ownerID := uuid.New()

	cmd := events.NewOperations(repo, nil, nil)

	ctx := idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: ownerID},
	})

	event := &models.Event{
		ID:      eventID,
		OwnerID: ownerID,
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)
	mock.When(repo.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleOwner, nil)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil idx client, but didn't panic")
		}
	}()

	_ = cmd.RemoveMember(ctx, eventID, models.RemoveMemberRequest{
		Email: "someone@example.com",
	})
}

func TestRemoveMember_NoIdentity(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()

	cmd := events.NewOperations(repo, nil, nil)

	err := cmd.RemoveMember(context.Background(), uuid.New(), models.RemoveMemberRequest{
		Email: "someone@example.com",
	})
	if err == nil {
		t.Fatal("expected identity error, got nil")
	}
}
