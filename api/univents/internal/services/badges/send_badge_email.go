package badges

import (
	"context"
	"encoding/base64"
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
// (profile URL, base64 data-URI) as an offline copy plus a link to the badge.
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
	qrDataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrPNG)

	body, err := assets.RenderBadgeEmittedEmail(assets.BadgeEmittedEmailData{
		AttendeeName: reg.AttendeeName,
		EventName:    event.FullName,
		EditionName:  edition.Name,
		BadgeLink:    badgeLink,
		QRDataURI:    qrDataURI,
	})
	if err != nil {
		telemetry.Log().Error("failed to render badge email", zap.Error(err))
		return
	}

	err = o.email.Send(email.Message{
		To:      []string{reg.AttendeeEmail},
		Subject: fmt.Sprintf("Your badge for %s — %s", event.FullName, edition.Name),
		Body:    body,
		HTML:    true,
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
