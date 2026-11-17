package jobs_test

import (
	"context"
	"errors"
	"testing"
	"time"
	"uuid"

	"github.com/ovechkin-dm/mockio/mock"

	"lib/email"

	"univents/internal/features/badges/jobs"
	"univents/models"
	"univents/ports"
)

type fakeSender struct {
	sent []email.Message
	err  error
}

func (f *fakeSender) SendWithInlineImage(msg email.Message, _ ...email.InlineImage) error {
	if f.err != nil {
		return f.err
	}
	f.sent = append(f.sent, msg)
	return nil
}

type fixture struct {
	emissions     ports.BadgeEmissionRepo
	registrations ports.RegistrationRepo
	sender        *fakeSender
	worker        *jobs.SendBadgeEmailWorker
	args          jobs.SendBadgeEmailArgs
}

func setup(t *testing.T, emission *models.BadgeEmission, reg *models.Registration) fixture {
	t.Helper()
	mock.SetUp(t)
	f := fixture{
		emissions:     mock.Mock[ports.BadgeEmissionRepo](),
		registrations: mock.Mock[ports.RegistrationRepo](),
		sender:        &fakeSender{},
		args:          jobs.SendBadgeEmailArgs{EmissionID: emission.ID, RegistrationID: reg.ID},
	}
	editions := mock.Mock[ports.EditionRepo]()
	events := mock.Mock[ports.EventRepo]()
	mock.When(f.emissions.GetEmission(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(emission, nil)
	mock.When(f.registrations.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(reg, nil)
	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: reg.EditionID, Name: "2026"}, nil)
	mock.When(events.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Event{FullName: "SCTI"}, nil)
	f.worker = jobs.NewSendBadgeEmailWorker(f.registrations, f.emissions, editions, events, f.sender, "https://events.test/")
	return f
}

func confirmed() (*models.BadgeEmission, *models.Registration) {
	userID := uuid.New()
	reg := &models.Registration{
		ID: uuid.New(), EditionID: uuid.New(), AttendeeUserID: &userID,
		AttendeeEmail: "attendee@example.com", AttendeeName: "Attendee",
		Status: models.RegistrationStatusConfirmed,
	}
	emission := &models.BadgeEmission{ID: uuid.New(), Status: models.BadgeEmissionStatusActive}
	return emission, reg
}

func TestSendBadgeEmail_SendsAndMarks(t *testing.T) {
	emission, reg := confirmed()
	f := setup(t, emission, reg)

	err := f.worker.Work(context.Background(), f.args)
	if err != nil {
		t.Fatalf("Work: %v", err)
	}
	if len(f.sender.sent) != 1 || f.sender.sent[0].To[0] != "attendee@example.com" {
		t.Fatalf("want one email to the attendee, got %+v", f.sender.sent)
	}
	_ = mock.Verify(f.emissions, mock.Once()).MarkEmailSent(mock.AnyContext(), mock.Equal(emission.ID))
}

func TestSendBadgeEmail_AlreadySentIsNoOp(t *testing.T) {
	emission, reg := confirmed()
	sentAt := time.Now()
	emission.EmailSentAt = &sentAt
	f := setup(t, emission, reg)

	err := f.worker.Work(context.Background(), f.args)
	if err != nil {
		t.Fatalf("Work: %v", err)
	}
	if len(f.sender.sent) != 0 {
		t.Fatalf("want no email, got %+v", f.sender.sent)
	}
}

func TestSendBadgeEmail_UnconfirmedIsNoOp(t *testing.T) {
	emission, reg := confirmed()
	reg.Status = models.RegistrationStatusPending
	f := setup(t, emission, reg)

	err := f.worker.Work(context.Background(), f.args)
	if err != nil {
		t.Fatalf("Work: %v", err)
	}
	if len(f.sender.sent) != 0 {
		t.Fatalf("want no email, got %+v", f.sender.sent)
	}
}

func TestSendBadgeEmail_SendErrorIsRetried(t *testing.T) {
	emission, reg := confirmed()
	f := setup(t, emission, reg)
	f.sender.err = errors.New("smtp down")

	err := f.worker.Work(context.Background(), f.args)
	if err == nil {
		t.Fatal("want the send error returned so the carrier retries")
	}
	_ = mock.Verify(f.emissions, mock.Never()).MarkEmailSent(mock.AnyContext(), mock.Any[uuid.UUID]())
}
