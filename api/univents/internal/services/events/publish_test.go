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

func TestPublish_OwnerCanPublishDraft(t *testing.T) {
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
		Status:  models.EventStatusDraft,
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)
	mock.When(repo.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleOwner, nil)
	mock.When(repo.Publish(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)

	err := cmd.Publish(ctx, eventID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_ = mock.Verify(repo, mock.Once()).Publish(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestPublish_AdminCanPublishDraft(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()
	authzSvc := authz.New(repo)

	eventID := uuid.New()
	ownerID := uuid.New()
	adminID := uuid.New()

	cmd := events.NewOperations(repo, nil, nil, authzSvc, mock.Mock[ports.BadgeStaffOps]())

	ctx := idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: adminID},
	})

	event := &models.Event{
		ID:      eventID,
		OwnerID: ownerID,
		Status:  models.EventStatusDraft,
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)
	mock.When(repo.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleAdmin, nil)
	mock.When(repo.Publish(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(nil)

	err := cmd.Publish(ctx, eventID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_ = mock.Verify(repo, mock.Once()).Publish(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestPublish_NonAdminForbidden(t *testing.T) {
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
		Status:  models.EventStatusDraft,
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)
	mock.When(repo.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleStaff, nil)

	err := cmd.Publish(ctx, eventID)
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}

	_ = mock.Verify(repo, mock.Never()).Publish(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestPublish_CannotPublishNonDraft(t *testing.T) {
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
		Status:  models.EventStatusActive,
	}

	mock.When(repo.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(event, nil)

	err := cmd.Publish(ctx, eventID)
	if err == nil {
		t.Fatal("expected bad request error, got nil")
	}

	_ = mock.Verify(repo, mock.Never()).Publish(mock.AnyContext(), mock.Any[uuid.UUID]())
}

func TestPublish_NoIdentityInCtx(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()
	authzSvc := authz.New(repo)

	cmd := events.NewOperations(repo, nil, nil, authzSvc, mock.Mock[ports.BadgeStaffOps]())
	ctx := context.Background()

	err := cmd.Publish(ctx, uuid.New())
	if err == nil {
		t.Fatal("expected identity error, got nil")
	}
}
