package events_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/internal/services/events"
	"univents/models"
	"univents/ports"
)

func TestAddMember_StaffForbidden(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()
	authzSvc := authz.New(repo)

	eventID := uuid.New()
	ownerID := uuid.New()
	staffID := uuid.New()

	cmd := events.NewOperations(repo, nil, nil, authzSvc, mock.Mock[ports.BadgeStaffOps]())

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

	_, err := cmd.AddMember(ctx, eventID, models.AddEventMemberInput{
		Email: "someone@example.com",
		Role:  models.EventMemberRoleStaff,
	})
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
}

func TestAddMember_OwnerGetsPastAuth(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()
	authzSvc := authz.New(repo)

	eventID := uuid.New()
	ownerID := uuid.New()

	cmd := events.NewOperations(repo, nil, nil, authzSvc, mock.Mock[ports.BadgeStaffOps]())

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

	// Owner passes CheckEvent via the owner role, then it tries
	// c.idx.Actors.GetByEmail which is nil — but that's fine,
	// the test proves authorization logic lets the owner through.
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on nil idx client, but didn't panic")
		}
	}()

	_, _ = cmd.AddMember(ctx, eventID, models.AddEventMemberInput{
		Email: "someone@example.com",
		Role:  models.EventMemberRoleStaff,
	})
}

func TestAddMember_NoIdentity(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()
	authzSvc := authz.New(repo)

	cmd := events.NewOperations(repo, nil, nil, authzSvc, mock.Mock[ports.BadgeStaffOps]())

	_, err := cmd.AddMember(context.Background(), uuid.New(), models.AddEventMemberInput{
		Email: "someone@example.com",
		Role:  models.EventMemberRoleStaff,
	})
	if err == nil {
		t.Fatal("expected identity error, got nil")
	}
}
