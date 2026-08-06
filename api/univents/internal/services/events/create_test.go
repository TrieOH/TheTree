package events_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"lib/database"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/internal/services/events"
	"univents/models"
	"univents/ports"
)

func TestCreate_Success(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()
	authzSvc := authz.New(repo)
	var txr = mock.Mock[database.TxRunner]()

	database.SetDefaultRunner(txr)

	ownerID := uuid.New()
	eventID := uuid.New()

	cmd := events.NewOperations(repo, nil, nil, authzSvc, mock.Mock[ports.BadgeStaffOps]())

	ctx := idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: ownerID},
	})

	payload := models.CreateEventInput{
		FullName: "Test Event",
		Slug:     "test-event",
	}

	// Execute the tx closure so repo calls actually happen.
	mock.When(txr.WithinTx(mock.AnyContext(), mock.Any[func(context.Context) error]())).
		ThenAnswer(func(args []any) []any {
			fn := args[1].(func(context.Context) error)
			return []any{fn(args[0].(context.Context))}
		})

	mock.When(repo.Create(mock.AnyContext(), mock.Any[*models.Event]())).
		ThenReturn(&models.Event{ID: eventID, OwnerID: ownerID}, nil)

	mock.When(repo.AddEventMember(
		mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID](), mock.Any[models.EventMemberRole](),
	)).ThenReturn(&models.EventMember{}, nil)

	event, err := cmd.Create(ctx, payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event == nil {
		t.Fatal("expected event, got nil")
	}
	if event.ID != eventID {
		t.Fatalf("expected event ID %v, got %v", eventID, event.ID)
	}

	_, _ = mock.Verify(repo, mock.Once()).Create(mock.AnyContext(), mock.Any[*models.Event]())
	_, _ = mock.Verify(repo, mock.Once()).AddEventMember(
		mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID](), mock.Any[models.EventMemberRole](),
	)
}

func TestCreate_NoIdentity(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()
	authzSvc := authz.New(repo)

	cmd := events.NewOperations(repo, nil, nil, authzSvc, mock.Mock[ports.BadgeStaffOps]())

	_, err := cmd.Create(context.Background(), models.CreateEventInput{
		FullName: "Test Event",
		Slug:     "test-event",
	})
	if err == nil {
		t.Fatal("expected identity error, got nil")
	}
}

func TestCreate_RepoCreateFails(t *testing.T) {
	mock.SetUp(t)

	var repo = mock.Mock[ports.EventRepo]()
	authzSvc := authz.New(repo)
	var txr = mock.Mock[database.TxRunner]()

	database.SetDefaultRunner(txr)

	ownerID := uuid.New()

	cmd := events.NewOperations(repo, nil, nil, authzSvc, mock.Mock[ports.BadgeStaffOps]())

	ctx := idx.WithIdentity(context.Background(), &idx.Identity{
		Sub: idx.Subject{ID: ownerID},
	})

	payload := models.CreateEventInput{
		FullName: "Test Event",
		Slug:     "test-event",
	}

	mock.When(txr.WithinTx(mock.AnyContext(), mock.Any[func(context.Context) error]())).
		ThenAnswer(func(args []any) []any {
			fn := args[1].(func(context.Context) error)
			return []any{fn(args[0].(context.Context))}
		})

	mock.When(repo.Create(mock.AnyContext(), mock.Any[*models.Event]())).
		ThenReturn(nil, errAssertAnError)

	_, err := cmd.Create(ctx, payload)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

var errAssertAnError = &testError{}

type testError struct{}

func (e *testError) Error() string { return "test error" }
