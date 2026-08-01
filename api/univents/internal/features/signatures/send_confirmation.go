package signatures

import (
	"context"
	"fmt"
	"os"
	"time"
	"univents/assets"
	"univents/models"

	"lib/crypto"
	"lib/email"
	"lib/telemetry"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func (o *Operations) sendConfirmationEmail(ctx context.Context, sig *models.Signature) {
	if sig.SignatoryEmail == nil {
		return
	}

	edition, err := o.editions.GetByID(ctx, sig.EditionID)
	if err != nil {
		return
	}

	event, err := o.events.GetByID(ctx, edition.EventID)
	if err != nil {
		return
	}

	token, err := crypto.SignHMACJWT(models.SignatureRevocationClaims{
		SignatureID: sig.ID,
		EditionID:   sig.EditionID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, []byte(o.hmacSecret))
	if err != nil {
		return
	}

	revokeLink := fmt.Sprintf("%s/signatures/revoke?token=%s", os.Getenv("APP_URL"), token)
	body, err := assets.RenderSignatureCreatedEmail(assets.SignatureCreatedEmailData{
		SignatoryName: sig.SignatoryName,
		EventName:     event.FullName,
		EditionName:   edition.Name,
		RevokeLink:    revokeLink,
	})
	if err != nil {
		return
	}

	err = o.email.Send(email.Message{
		To:      []string{*sig.SignatoryEmail},
		Subject: fmt.Sprintf("Your signature for %s — %s", event.FullName, edition.Name),
		Body:    body,
		HTML:    true,
	})
	if err != nil {
		telemetry.Log().Error("failed to send signature confirmation email", zap.Error(err))
	}
}
