package jobs

import (
	"context"
	"fmt"
	"os"
	"strings"

	"lib/email"
	"lib/telemetry"

	"univents/assets"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// EmailSender is the SMTP send seam. Satisfied by *email.Client; faked in
// tests.
type EmailSender interface {
	Send(msg email.Message) error
}

// SendGiftEmailArgs is the `gifts.send_email` job: email the accountless
// recipient of a paid gift, telling them they were gifted a ticket and to
// create an account under the gifted email to claim it (the claim ties the
// ticket to their account and emits the deferred badge). Enqueued in the
// same tx that confirms the registration (webhook approval / free order),
// so the email never goes out for an abandoned reservation.
type SendGiftEmailArgs struct {
	RegistrationID uuid.UUID `json:"registration_id"`
}

func (SendGiftEmailArgs) Kind() string { return "gifts.send_email" }

// SendGiftEmailWorker sends the gifted-ticket email. Skips registrations
// whose attendee already has an account (claim beat the job — they'll get
// the badge email instead) or that are no longer confirmed. Send errors are
// returned so River retries with backoff.
type SendGiftEmailWorker struct {
	river.WorkerDefaults[SendGiftEmailArgs]

	registrations ports.RegistrationRepo
	editions      ports.EditionRepo
	events        ports.EventRepo
	ticketTypes   ports.TicketTypeRepo
	email         EmailSender
}

func NewSendGiftEmailWorker(
	registrations ports.RegistrationRepo,
	editions ports.EditionRepo,
	events ports.EventRepo,
	ticketTypes ports.TicketTypeRepo,
	email EmailSender,
) *SendGiftEmailWorker {
	return &SendGiftEmailWorker{
		registrations: registrations,
		editions:      editions,
		events:        events,
		ticketTypes:   ticketTypes,
		email:         email,
	}
}

// Work sends the email. Idempotent in effect: once the recipient claims
// (attendee_user_id set), the registration no longer matches and a retry is
// a no-op.
func (w *SendGiftEmailWorker) Work(ctx context.Context, job *river.Job[SendGiftEmailArgs]) error {
	ctx, span := telemetry.StartSpan(ctx, "SendGiftEmailWorker.Work")
	defer span.End()

	reg, err := w.registrations.GetByID(ctx, job.Args.RegistrationID)
	if err != nil {
		return err // missing registration is genuinely broken — let River retry
	}
	if reg.AttendeeUserID != nil {
		// Claimed before the job ran — they have an account now; the badge
		// email covers the confirmation. No gift email needed.
		telemetry.Log().Info("gift email skipped: recipient already has an account",
			zap.String("registration_id", reg.ID.String()))
		return nil
	}
	if reg.Status != models.RegistrationStatusConfirmed {
		// Not confirmed (shouldn't happen — enqueued on confirm), nothing to
		// announce; a refund/expiry would make this a stale announce.
		return nil
	}

	edition, err := w.editions.GetByID(ctx, reg.EditionID)
	if err != nil {
		return err
	}
	event, err := w.events.GetByID(ctx, edition.EventID)
	if err != nil {
		return err
	}
	ticketType, err := w.ticketTypes.GetByID(ctx, reg.TicketTypeID)
	if err != nil {
		return err
	}

	body, err := assets.RenderTicketGiftedEmail(assets.TicketGiftedEmailData{
		RecipientName:  reg.AttendeeName,
		RecipientEmail: reg.AttendeeEmail,
		EventName:      event.FullName,
		EditionName:    edition.Name,
		TicketTypeName: ticketType.Name,
		// The claim is lazy: the event page's my-ticket read ties the
		// ticket to the logged-in account and emits the deferred badge.
		ClaimLink: strings.TrimRight(os.Getenv("APP_URL"), "/") + "/events/" + event.Slug,
	})
	if err != nil {
		return err
	}

	err = w.email.Send(email.Message{
		To:      []string{reg.AttendeeEmail},
		Subject: fmt.Sprintf("You received a gifted ticket for %s — %s", event.FullName, edition.Name),
		Body:    body,
		HTML:    true,
	})
	if err != nil {
		telemetry.Log().Error("failed to send gifted-ticket email",
			zap.String("registration_id", reg.ID.String()),
			zap.String("to", reg.AttendeeEmail),
			zap.Error(err))
		return err
	}
	return nil
}
