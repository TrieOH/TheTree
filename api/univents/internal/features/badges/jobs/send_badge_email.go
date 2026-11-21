// Package jobs holds the badges feature's background work.
package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"uuid"

	"github.com/skip2/go-qrcode"

	"lib/database"
	"lib/email"
	"lib/telemetry"

	"univents/assets"
	"univents/models"
	"univents/ports"
)

// EmailSender is the send seam. Satisfied by *email.Client; faked in tests.
type EmailSender interface {
	SendWithInlineImage(ctx context.Context, msg email.Message, images ...email.InlineImage) error
}

// SendBadgeEmailArgs is the `badges.send_email` job: email the participant
// their badge. Enqueued in the same tx that emits the badge (checkout,
// webhook approval, gift claim), so it never goes out for a rolled-back
// confirmation.
type SendBadgeEmailArgs struct {
	EmissionID     uuid.UUID `json:"emission_id"`
	RegistrationID uuid.UUID `json:"registration_id"`
}

func (SendBadgeEmailArgs) Kind() string { return "badges.send_email" }

// SendBadgeEmailWorker sends the badge email exactly once. Send errors are
// returned so the job carrier retries with backoff.
type SendBadgeEmailWorker struct {
	registrations ports.RegistrationRepo
	emissions     ports.BadgeEmissionRepo
	editions      ports.EditionRepo
	events        ports.EventRepo
	email         EmailSender
	appURL        string
}

func NewSendBadgeEmailWorker(
	registrations ports.RegistrationRepo,
	emissions ports.BadgeEmissionRepo,
	editions ports.EditionRepo,
	events ports.EventRepo,
	email EmailSender,
	appURL string,
) *SendBadgeEmailWorker {
	return &SendBadgeEmailWorker{
		registrations: registrations,
		emissions:     emissions,
		editions:      editions,
		events:        events,
		email:         email,
		appURL:        appURL,
	}
}

// Work sends the email: the QR (profile URL) as an inline image (cid: —
// email clients render these, unlike data: URIs) plus a link to the badge.
// A missing, revoked or already-emailed emission, or a registration that
// is no longer confirmed, is a no-op.
func (w *SendBadgeEmailWorker) Work(ctx context.Context, args SendBadgeEmailArgs) error {
	emission, err := w.emissions.GetEmission(ctx, args.EmissionID)
	if errors.Is(err, database.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if emission.EmailSentAt != nil || emission.Status == models.BadgeEmissionStatusRevoked {
		return nil
	}

	reg, err := w.registrations.GetByID(ctx, args.RegistrationID)
	if errors.Is(err, database.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if reg.Status != models.RegistrationStatusConfirmed || reg.AttendeeUserID == nil {
		telemetry.Log().InfoContext(ctx, "badge email skipped: registration not confirmed or attendee has no account",
			slog.String("registration_id", reg.ID.String()))
		return nil
	}

	edition, err := w.editions.GetByID(ctx, reg.EditionID)
	if err != nil {
		return fmt.Errorf("badge email: get edition: %w", err)
	}
	event, err := w.events.GetByID(ctx, edition.EventID)
	if err != nil {
		return fmt.Errorf("badge email: get event: %w", err)
	}

	actionURL := strings.TrimRight(w.appURL, "/") + "/profile/" + reg.AttendeeUserID.String()
	qrPNG, err := qrcode.Encode(actionURL, qrcode.Medium, 256)
	if err != nil {
		return fmt.Errorf("badge email: render QR: %w", err)
	}

	body, err := assets.RenderBadgeEmittedEmail(assets.BadgeEmittedEmailData{
		AttendeeName: reg.AttendeeName,
		EventName:    event.FullName,
		EditionName:  edition.Name,
		BadgeLink:    actionURL + "/badges/" + emission.ID.String(),
	})
	if err != nil {
		return fmt.Errorf("badge email: render: %w", err)
	}

	err = w.email.SendWithInlineImage(ctx, email.Message{
		To:      []string{reg.AttendeeEmail},
		Subject: fmt.Sprintf("Your badge for %s — %s", event.FullName, edition.Name),
		Body:    body,
		HTML:    true,
	}, email.InlineImage{
		ContentID: "badge-qr",
		MIMEType:  "image/png",
		Data:      qrPNG,
	})
	if err != nil {
		return fmt.Errorf("badge email: send: %w", err)
	}

	return w.emissions.MarkEmailSent(ctx, emission.ID)
}
