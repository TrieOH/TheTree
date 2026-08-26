package badges_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"lib/email"

	"univents/models"
	"univents/ports"
)

func TestEmitForConfirmedRegistration_UpsertsAndEmailsOnce(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	userID := uuid.New()
	regID := uuid.New()

	reg := &models.Registration{
		ID:             regID,
		EditionID:      editionID,
		AttendeeUserID: &userID,
		AttendeeEmail:  "attendee@example.com",
		AttendeeName:   "Attendee",
		Status:         models.RegistrationStatusConfirmed,
	}

	registrations := mock.Mock[ports.RegistrationRepo]()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()

	mock.When(registrations.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(reg, nil)

	// The email goroutine needs real editions/events lookups and a live client
	// (the SMTP dial fails harmlessly and is logged).
	editions := mock.Mock[ports.EditionRepo]()
	events := mock.Mock[ports.EventRepo]()
	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: editionID, Name: "Edition"}, nil)
	mock.When(events.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Event{FullName: "Event"}, nil)
	emailClient := email.NewClient(email.Config{Host: "127.0.0.1", Port: 1, From: "badges@test.local"})

	var upserted *models.BadgeEmission
	mock.When(emissions.Upsert(mock.AnyContext(), mock.Any[*models.BadgeEmission]())).
		ThenAnswer(func(args []any) []any {
			e := args[1].(*models.BadgeEmission)
			upserted = e
			// First emission: no email yet.
			em := &models.BadgeEmission{ID: uuid.New(), EditionID: e.EditionID, UserID: e.UserID, Origin: e.Origin, RegistrationID: e.RegistrationID}
			return []any{em, nil}
		})

	ops := newOps(t, nil, emissions, registrations, editions, events, emailClient)

	emission, err := ops.EmitForConfirmedRegistration(context.Background(), regID)
	if err != nil {
		t.Fatalf("EmitForConfirmedRegistration: %v", err)
	}
	if emission == nil {
		t.Fatal("want emission, got nil")
	}
	if upserted == nil {
		t.Fatal("upsert not called")
	}
	if upserted.EditionID != editionID || upserted.UserID != userID || upserted.Origin != models.BadgeEmissionOriginParticipant {
		t.Errorf("unexpected upsert args: %+v", upserted)
	}
	if upserted.RegistrationID == nil || *upserted.RegistrationID != regID {
		t.Errorf("want registration_id %s, got %v", regID, upserted.RegistrationID)
	}
}

func TestEmitForConfirmedRegistration_PendingIsNoOp(t *testing.T) {
	mock.SetUp(t)

	regID := uuid.New()
	reg := &models.Registration{ID: regID, Status: models.RegistrationStatusPending}

	registrations := mock.Mock[ports.RegistrationRepo]()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	mock.When(registrations.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(reg, nil)

	ops := newOps(t, nil, emissions, registrations, nil, nil, nil)

	emission, err := ops.EmitForConfirmedRegistration(context.Background(), regID)
	if err != nil {
		t.Fatalf("EmitForConfirmedRegistration: %v", err)
	}
	if emission != nil {
		t.Fatalf("pending registration must not emit, got %+v", emission)
	}
}

// TestEmitForConfirmedRegistration_AccountlessSkips pins the email-only
// gift deferral: a confirmed registration whose attendee has no account yet
// (nil attendee_user_id) never emits — there is no profile to attach the
// badge to or email it to. Emission resumes once the recipient claims an
// account (the claim flow re-runs this from a registration that has one).
func TestEmitForConfirmedRegistration_AccountlessSkips(t *testing.T) {
	mock.SetUp(t)

	regID := uuid.New()
	reg := &models.Registration{
		ID:        regID,
		EditionID: uuid.New(),
		Status:    models.RegistrationStatusConfirmed,
		// AttendeeUserID nil = accountless gifted recipient.
	}

	registrations := mock.Mock[ports.RegistrationRepo]()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	mock.When(registrations.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(reg, nil)
	// Upsert is stubbed to return a real emission: a non-nil return here
	// would prove the skip did not happen (emission must stay nil).
	mock.When(emissions.Upsert(mock.AnyContext(), mock.Any[*models.BadgeEmission]())).
		ThenReturn(&models.BadgeEmission{ID: uuid.New()}, nil)

	ops := newOps(t, nil, emissions, registrations, nil, nil, nil)

	emission, err := ops.EmitForConfirmedRegistration(context.Background(), regID)
	if err != nil {
		t.Fatalf("EmitForConfirmedRegistration: %v", err)
	}
	if emission != nil {
		t.Fatalf("accountless registration must not emit, got %+v", emission)
	}
}

func TestRevokeForRegistration(t *testing.T) {
	mock.SetUp(t)

	regID := uuid.New()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()

	var revokedID uuid.UUID
	var revokedReason string
	mock.When(emissions.RevokeByRegistration(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[string]())).
		ThenAnswer(func(args []any) []any {
			revokedID = args[1].(uuid.UUID)
			revokedReason = args[2].(string)
			return []any{nil}
		})

	ops := newOps(t, nil, emissions, nil, nil, nil, nil)

	err := ops.RevokeForRegistration(context.Background(), regID, "cancelled")
	if err != nil {
		t.Fatalf("RevokeForRegistration: %v", err)
	}
	if revokedID != regID || revokedReason != "cancelled" {
		t.Errorf("unexpected revoke args: %s %q", revokedID, revokedReason)
	}
}

func TestEmitForConfirmedRegistration_NoSecondEmail(t *testing.T) {
	mock.SetUp(t)

	regID := uuid.New()
	reg := &models.Registration{
		ID: regID, EditionID: uuid.New(), AttendeeUserID: new(uuid.New()),
		Status: models.RegistrationStatusConfirmed,
	}

	registrations := mock.Mock[ports.RegistrationRepo]()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()

	mock.When(registrations.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(reg, nil)

	alreadySent := time.Now()
	mock.When(emissions.Upsert(mock.AnyContext(), mock.Any[*models.BadgeEmission]())).
		ThenAnswer(func(args []any) []any {
			e := args[1].(*models.BadgeEmission)
			em := &models.BadgeEmission{
				ID: uuid.New(), EditionID: e.EditionID, UserID: e.UserID,
				Origin: e.Origin, RegistrationID: e.RegistrationID,
				EmailSentAt: &alreadySent,
			}
			return []any{em, nil}
		})

	ops := newOps(t, nil, emissions, registrations, nil, nil, nil)

	_, err := ops.EmitForConfirmedRegistration(context.Background(), regID)
	if err != nil {
		t.Fatalf("EmitForConfirmedRegistration: %v", err)
	}

	// The badge was already emailed once: no goroutine fires, no re-send.
	_ = mock.Verify(emissions, mock.Never()).MarkEmailSent(mock.AnyContext(), mock.Any[uuid.UUID]())
}
