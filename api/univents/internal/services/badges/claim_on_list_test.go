package badges_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"lib/email"

	idx "sdk/identityx"
	"univents/models"
	"univents/ports"
)

// TestListByUser_ClaimsGiftAndEmitsDeferredBadge pins the profile-badges gift
// claim: a confirmed email-only gift (attendee_user_id NULL, badge deferred
// at approval) is claimed from the account's own email when the recipient's
// profile loads, and the deferred badge is emitted into the response. The
// account email casing differs from the gift email to exercise the
// normalization (TrimSpace + ToLower), mirroring the my-ticket claim.
func TestListByUser_ClaimsGiftAndEmitsDeferredBadge(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	eventID := uuid.New()
	ticketTypeID := uuid.New()
	userID := uuid.New()
	regID := uuid.New()

	// The account exists (a profile requires one) with the gift email.
	actors := &fakeActors{byID: map[uuid.UUID]*idx.Actor{
		userID: {ID: userID, Email: new("Sophia@Example.com")},
	}}

	// Pre-claim state: confirmed gift, no account tied yet.
	claimed := &models.Registration{
		ID:             regID,
		EditionID:      editionID,
		TicketTypeID:   ticketTypeID,
		AttendeeUserID: &userID,
		AttendeeEmail:  "sophia@example.com",
		AttendeeName:   "Sophia",
		Status:         models.RegistrationStatusConfirmed,
	}

	registrations := mock.Mock[ports.RegistrationRepo]()
	var claimedEmail string
	var claimedUser uuid.UUID
	mock.When(registrations.ClaimAllByAttendeeEmail(mock.AnyContext(), mock.Any[string](), mock.Any[uuid.UUID]())).
		ThenAnswer(func(args []any) []any {
			claimedEmail = args[1].(string)
			claimedUser = args[2].(uuid.UUID)
			return []any{[]*models.Registration{claimed}, nil}
		})
	mock.When(registrations.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(claimed, nil)

	emissionID := uuid.New()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	mock.When(emissions.Upsert(mock.AnyContext(), mock.Any[*models.BadgeEmission]())).
		ThenReturn(&models.BadgeEmission{
			ID:             emissionID,
			EditionID:      editionID,
			UserID:         userID,
			Origin:         models.BadgeEmissionOriginParticipant,
			RegistrationID: &regID,
		}, nil)
	// After the claim the read returns the fresh emission.
	mock.When(emissions.ListViewsByUser(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{
			{
				BadgeEmission: models.BadgeEmission{
					ID:             emissionID,
					EditionID:      editionID,
					UserID:         userID,
					Origin:         models.BadgeEmissionOriginParticipant,
					RegistrationID: &regID,
					Status:         models.BadgeEmissionStatusActive,
				},
				EditionName:  "2026",
				EndsAt:       time.Now().Add(24 * time.Hour),
				EventName:    "SCTI",
				TicketTypeID: &ticketTypeID,
				TicketName:   new("Básico"),
			},
		}, nil)

	// The badge email goroutine needs editions/events lookups and a live
	// client (the SMTP dial fails harmlessly and is logged).
	editions := mock.Mock[ports.EditionRepo]()
	events := mock.Mock[ports.EventRepo]()
	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: editionID, Name: "2026", EventID: eventID}, nil)
	mock.When(events.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Event{FullName: "SCTI"}, nil)
	emailClient := email.NewClient(email.Config{Host: "127.0.0.1", Port: 1, From: "badges@test.local"})

	// The profile read resolves templates for the rendered editions.
	templates := mock.Mock[ports.BadgeTemplateRepo]()
	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn([]models.BadgeTemplate{
			{ID: uuid.New(), EditionID: editionID, Name: "Básico", DesignData: []byte(`{}`)},
		}, nil)

	ops := newOpsWithDeps(t, templates, emissions, registrations, editions, events, actors, emailClient, mock.Mock[ports.EventRepo]())

	groups, err := ops.ListByUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListByUser: %v", err)
	}
	if len(groups.Attendant.Current) != 1 {
		t.Fatalf("want 1 current attendant badge, got %d", len(groups.Attendant.Current))
	}
	badge := groups.Attendant.Current[0]
	if badge.EmissionID != emissionID {
		t.Errorf("want emission %s, got %s", emissionID, badge.EmissionID)
	}
	if badge.TicketName == nil || *badge.TicketName != "Básico" {
		t.Errorf("want ticket Básico, got %v", badge.TicketName)
	}

	// The claim ran with the normalized account email, tied to the user.
	if claimedEmail != "sophia@example.com" {
		t.Errorf("want claim email normalized to sophia@example.com, got %q", claimedEmail)
	}
	if claimedUser != userID {
		t.Errorf("want claim user %s, got %s", userID, claimedUser)
	}
	_, _ = mock.Verify(registrations, mock.Times(1)).
		ClaimAllByAttendeeEmail(mock.AnyContext(), mock.Any[string](), mock.Any[uuid.UUID]())
	_, _ = mock.Verify(emissions, mock.Times(1)).Upsert(mock.AnyContext(), mock.Any[*models.BadgeEmission]())
}

// TestListByUser_ClaimsPendingGiftWithoutEmitting pins that a pending gift is
// claimed (tied to the account) but emits no badge — emission waits for
// confirmation, exactly like EmitForConfirmedRegistration's no-op contract.
func TestListByUser_ClaimsPendingGiftWithoutEmitting(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	userID := uuid.New()
	regID := uuid.New()

	actors := &fakeActors{byID: map[uuid.UUID]*idx.Actor{
		userID: {ID: userID, Email: new("sophia@example.com")},
	}}

	pending := &models.Registration{
		ID:             regID,
		EditionID:      editionID,
		AttendeeUserID: &userID,
		AttendeeEmail:  "sophia@example.com",
		Status:         models.RegistrationStatusPending,
	}

	registrations := mock.Mock[ports.RegistrationRepo]()
	var claimedEmail string
	mock.When(registrations.ClaimAllByAttendeeEmail(mock.AnyContext(), mock.Any[string](), mock.Any[uuid.UUID]())).
		ThenAnswer(func(args []any) []any {
			claimedEmail = args[1].(string)
			return []any{[]*models.Registration{pending}, nil}
		})
	mock.When(registrations.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(pending, nil)

	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	mock.When(emissions.ListViewsByUser(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{}, nil)

	ops := newOpsWithDeps(t, nil, emissions, registrations, nil, nil, actors, nil, mock.Mock[ports.EventRepo]())

	groups, err := ops.ListByUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListByUser: %v", err)
	}
	if len(groups.Attendant.Current) != 0 || len(groups.Attendant.Past) != 0 {
		t.Fatalf("pending gift must not render a badge, got %+v", groups.Attendant)
	}

	if claimedEmail != "sophia@example.com" {
		t.Errorf("want claim email normalized to sophia@example.com, got %q", claimedEmail)
	}
	_, _ = mock.Verify(registrations, mock.Times(1)).
		ClaimAllByAttendeeEmail(mock.AnyContext(), mock.Any[string](), mock.Any[uuid.UUID]())
	_, _ = mock.Verify(emissions, mock.Never()).Upsert(mock.AnyContext(), mock.Any[*models.BadgeEmission]())
}
