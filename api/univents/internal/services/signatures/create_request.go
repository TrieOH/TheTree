package signatures

import (
	"context"
	"fmt"
	"lib/crypto"
	"lib/email"
	"lib/telemetry"
	"os"
	idx "sdk/identityx"
	"time"
	"univents/assets"
	"univents/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (o *Operations) CreateRequest(ctx context.Context, payload models.CreateSignatureRequestInput) (*models.SignatureRequest, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.CreateRequest")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	idempotencyKey := uuid.New()
	expiresAt := time.Now().Add(time.Duration(payload.ExpiresInDays) * 24 * time.Hour)
	request, err := o.requests.CreateRequest(ctx, &models.SignatureRequest{
		EditionID:       payload.EditionID,
		CreatedBy:       ident.Sub.ID,
		SignatoryName:   payload.SignatoryName,
		SignatoryTitle:  payload.SignatoryTitle,
		SignatoryEmail:  payload.SignatoryEmail,
		SignatoryUserID: payload.SignatoryUserID,
		IdempotencyKey:  idempotencyKey,
		ExpiresAt:       expiresAt,
	})
	if err != nil {
		return nil, err
	}

	event, err := o.events.GetByID(ctx, edition.EventID)
	if err != nil {
		return nil, err
	}

	token, err := crypto.SignHMACJWT(models.SignatureRequestClaims{
		RequestID: request.ID,
		EditionID: payload.EditionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}, []byte(o.hmacSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to encode signature token: %w", err)
	}

	link := fmt.Sprintf("%s/signature-requests/fulfill?token=%s", os.Getenv("APP_URL"), token)
	body, err := assets.RenderRequestSignatureEmail(assets.RequestSignatureEmailData{
		SignatoryName: payload.SignatoryName,
		EventName:     event.FullName,
		EditionName:   edition.Name,
		Link:          link,
		ExpiresInDays: payload.ExpiresInDays,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render email: %w", err)
	}

	err = o.email.Send(email.Message{
		To:      []string{*payload.SignatoryEmail},
		Subject: fmt.Sprintf("Signature request for %s — %s", event.FullName, edition.Name),
		Body:    body,
		HTML:    true,
	})
	if err != nil {
		telemetry.Log().Error("failed to send signature request email", zap.Error(err))
	}

	return request, nil
}
