package service

import (
	"context"
	"fmt"
	"lib/email"
	"lib/telemetry"
	"log/slog"
	"os"
	"strings"
	"univents/assets"
	"univents/models"
	"uuid"

	"github.com/skip2/go-qrcode"
)

// sendBadgeEmail emails the participant their badge exactly once: the QR
// (profile URL) as an inline image (cid: — email clients render these,
// unlike data: URIs) plus a link to the badge.
func (o *Operations) sendBadgeEmail(ctx context.Context, reg *models.Registration, emission *models.BadgeEmission) {
	if reg.AttendeeUserID == nil {
		// Unreachable via EmitForConfirmedRegistration (accountless
		// registrations never emit), kept as a guard for future callers.
		telemetry.Log().Info("badge email skipped: attendee has no account yet",
			slog.String("registration_id", reg.ID.String()))
		return
	}

	edition, err := o.Editions.GetByID(ctx, reg.EditionID)
	if err != nil {
		telemetry.Log().Error("failed to get edition for badge email", telemetry.Err(err))
		return
	}

	event, err := o.Events.GetByID(ctx, edition.EventID)
	if err != nil {
		telemetry.Log().Error("failed to get event for badge email", telemetry.Err(err))
		return
	}

	actionURL := profileURL(*reg.AttendeeUserID)
	badgeLink := actionURL + "/badges/" + emission.ID.String()

	qrPNG, err := qrcode.Encode(actionURL, qrcode.Medium, 256)
	if err != nil {
		telemetry.Log().Error("failed to render badge QR for email", telemetry.Err(err))
		return
	}

	body, err := assets.RenderBadgeEmittedEmail(assets.BadgeEmittedEmailData{
		AttendeeName: reg.AttendeeName,
		EventName:    event.FullName,
		EditionName:  edition.Name,
		BadgeLink:    badgeLink,
	})
	if err != nil {
		telemetry.Log().Error("failed to render badge email", telemetry.Err(err))
		return
	}

	err = o.Email.SendWithInlineImage(email.Message{
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
		telemetry.Log().Error("failed to send badge email", telemetry.Err(err))
		return
	}

	_ = o.Emissions.MarkEmailSent(context.Background(), emission.ID)
}

func profileURL(userID uuid.UUID) string {
	return strings.TrimRight(os.Getenv("APP_URL"), "/") + "/profile/" + userID.String()
}
