package badges_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"univents/models"
	"univents/ports"
)

func TestAwardStaffBadgesForUser_OnlyPublishedCurrentOrFuture(t *testing.T) {
	mock.SetUp(t)

	eventID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	past := &models.Edition{ID: uuid.New(), EventID: eventID, IsDraft: false, EndsAt: now.Add(-time.Hour)}
	active := &models.Edition{ID: uuid.New(), EventID: eventID, IsDraft: false, EndsAt: now.Add(time.Hour)}
	future := &models.Edition{ID: uuid.New(), EventID: eventID, IsDraft: false, EndsAt: now.Add(72 * time.Hour)}

	editions := mock.Mock[ports.EditionRepo]()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()

	mock.When(editions.ListPublic(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.Edition{*past, *active, *future}, nil)

	upserted := map[uuid.UUID]bool{}
	mock.When(emissions.Upsert(mock.AnyContext(), mock.Any[*models.BadgeEmission]())).
		ThenAnswer(func(args []any) []any {
			e := args[1].(*models.BadgeEmission)
			upserted[e.EditionID] = true
			return []any{e, nil}
		})

	ops := newOps(t, nil, emissions, nil, editions, nil, nil)

	err := ops.AwardStaffBadgesForUser(context.Background(), eventID, userID)
	if err != nil {
		t.Fatalf("AwardStaffBadgesForUser: %v", err)
	}

	if len(upserted) != 2 {
		t.Fatalf("want badges for active+future only, got %d upserts", len(upserted))
	}
	for _, id := range []uuid.UUID{active.ID, future.ID} {
		if !upserted[id] {
			t.Errorf("missing upsert for edition %s", id)
		}
	}
	if upserted[past.ID] {
		t.Error("past edition must not be awarded")
	}
}

func TestRevokeStaffBadgesForUser_KeepsPast(t *testing.T) {
	mock.SetUp(t)

	eventID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	past := &models.Edition{ID: uuid.New(), EventID: eventID, IsDraft: false, EndsAt: now.Add(-time.Hour)}
	active := &models.Edition{ID: uuid.New(), EventID: eventID, IsDraft: false, EndsAt: now.Add(time.Hour)}

	editions := mock.Mock[ports.EditionRepo]()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()

	mock.When(editions.ListPublic(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.Edition{*past, *active}, nil)

	revoked := map[uuid.UUID]bool{}
	mock.When(emissions.Revoke(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID](), mock.Any[models.BadgeEmissionOrigin](), mock.Any[string]())).
		ThenAnswer(func(args []any) []any {
			editionID := args[1].(uuid.UUID)
			revoked[editionID] = true
			return []any{nil}
		})

	ops := newOps(t, nil, emissions, nil, editions, nil, nil)

	err := ops.RevokeStaffBadgesForUser(context.Background(), eventID, userID)
	if err != nil {
		t.Fatalf("RevokeStaffBadgesForUser: %v", err)
	}

	if !revoked[active.ID] {
		t.Error("current edition badge must be revoked")
	}
	if revoked[past.ID] {
		t.Error("past edition badge must be kept (keepsake)")
	}
}

func TestAwardStaffBadgesForEdition_AllMembers(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	eventID := uuid.New()
	memberA := uuid.New()
	memberB := uuid.New()

	editions := mock.Mock[ports.EditionRepo]()
	events := mock.Mock[ports.EventRepo]()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()

	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: editionID, EventID: eventID}, nil)
	mock.When(events.ListEventMembers(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.EventMember{
			{UserID: memberA}, {UserID: memberB},
		}, nil)

	upserted := map[uuid.UUID]bool{}
	mock.When(emissions.Upsert(mock.AnyContext(), mock.Any[*models.BadgeEmission]())).
		ThenAnswer(func(args []any) []any {
			e := args[1].(*models.BadgeEmission)
			upserted[e.UserID] = true
			return []any{e, nil}
		})

	ops := newOps(t, nil, emissions, nil, editions, events, nil)

	err := ops.AwardStaffBadgesForEdition(context.Background(), editionID)
	if err != nil {
		t.Fatalf("AwardStaffBadgesForEdition: %v", err)
	}

	if !upserted[memberA] || !upserted[memberB] {
		t.Errorf("want both members upserted, got %v", upserted)
	}
}
