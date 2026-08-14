package badges

import (
	"context"
	"fmt"
	"lib/email"
	"lib/telemetry"
	"os"
	"strings"
	"univents/assets"
	"univents/models"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
)

// sendBadgeEmail emails the participant their badge exactly once: the QR
// (profile URL) as an inline image (cid: — email clients render these,
// unlike data: URIs) plus a link to the badge.
func (o *Operations) sendBadgeEmail(ctx context.Context, reg *models.Registration, emission *models.BadgeEmission) {
	edition, err := o.editions.GetByID(ctx, reg.EditionID)
	if err != nil {
		telemetry.Log().Error("failed to get edition for badge email", zap.Error(err))
		return
	}

	event, err := o.events.GetByID(ctx, edition.EventID)
	if err != nil {
		telemetry.Log().Error("failed to get event for badge email", zap.Error(err))
		return
	}

	actionURL := profileURL(reg.AttendeeUserID)
	badgeLink := actionURL + "/badges/" + emission.ID.String()

	qrPNG, err := qrcode.Encode(actionURL, qrcode.Medium, 256)
	if err != nil {
		telemetry.Log().Error("failed to render badge QR for email", zap.Error(err))
		return
	}

	body, err := assets.RenderBadgeEmittedEmail(assets.BadgeEmittedEmailData{
		AttendeeName: reg.AttendeeName,
		EventName:    event.FullName,
		EditionName:  edition.Name,
		BadgeLink:    badgeLink,
	})
	if err != nil {
		telemetry.Log().Error("failed to render badge email", zap.Error(err))
		return
	}

	err = o.email.SendWithInlineImage(email.Message{
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
		telemetry.Log().Error("failed to send badge email", zap.Error(err))
		return
	}

	_ = o.emissions.MarkEmailSent(context.Background(), emission.ID)
}

func profileURL(userID uuid.UUID) string {
	return strings.TrimRight(os.Getenv("APP_URL"), "/") + "/profile/" + userID.String()
}
