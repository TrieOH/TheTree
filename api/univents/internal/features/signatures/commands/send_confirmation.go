package commands

import (
	"context"
	"fmt"
	"time"
	"univents/assets"
	"univents/models"

	"lib/crypto"
	"lib/email"
	"lib/telemetry"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func (c *Commands) sendConfirmationEmail(ctx context.Context, sig *models.Signature) {
	if sig.SignatoryEmail == nil {
		return
	}

	edition, err := c.editions.GetByID(ctx, sig.EditionID)
	if err != nil {
		return
	}

	event, err := c.events.GetByID(ctx, edition.EventID)
	if err != nil {
		return
	}

	token, err := crypto.SignHMACJWT(models.SignatureRevocationClaims{
		SignatureID: sig.ID,
		EditionID:   sig.EditionID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}, []byte(c.hmacSecret))
	if err != nil {
		return
	}

	revokeLink := "https://yourapp.com/signatures/revoke?token=" + token
	body, err := assets.RenderSignatureCreatedEmail(assets.SignatureCreatedEmailData{
		SignatoryName: sig.SignatoryName,
		EventName:     event.FullName,
		EditionName:   edition.Name,
		RevokeLink:    revokeLink,
	})
	if err != nil {
		return
	}

	err = c.email.Send(email.Message{
		To:      []string{*sig.SignatoryEmail},
		Subject: fmt.Sprintf("Your signature for %s — %s", event.FullName, edition.Name),
		Body:    body,
		HTML:    true,
	})
	if err != nil {
		telemetry.Log().Error("failed to send signature confirmation email", zap.Error(err))
	}
}
