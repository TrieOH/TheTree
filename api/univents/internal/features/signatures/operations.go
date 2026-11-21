// Package signatures is the signatures feature: its repo, handlers and
// operations. Build wires it.
package signatures

import (
	"context"
	"fmt"
	"lib/crypto"
	"lib/email"
	"lib/telemetry"
	"time"
	"univents/assets"
	"univents/errs"
	"univents/internal/authz"
	"univents/models"
	"univents/ports"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

type Deps struct {
	Events     ports.EventRepo
	Editions   ports.EditionRepo
	Signatures ports.SignatureRepo
	Requests   ports.SignatureRequestRepo
	Email      *email.Client
	HmacSecret string
	AppURL     string // the Univents front, for the links in emails
	Authz      *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) CancelRequest(ctx context.Context, requestID uuid.UUID, reason *string) error {
	request, err := o.Requests.GetRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, request.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	if request.Status != models.SignatureRequestStatusPending {
		return errs.ErrSignatureRequestNotPending
	}

	return o.Requests.CancelRequest(ctx, requestID, reason)
}

func (o *Operations) Create(ctx context.Context, payload models.AddSignatureInput) (*models.Signature, error) {
	edition, err := o.Editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	ident, err := o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	sig, err := o.Signatures.Create(ctx, &models.Signature{
		EditionID:       payload.EditionID,
		CreatedBy:       ident.Sub.ID,
		SignatoryName:   payload.SignatoryName,
		SignatoryTitle:  payload.SignatoryTitle,
		SignatoryEmail:  payload.SignatoryEmail,
		SignatoryUserID: payload.SignatoryUserID,
		ImageURL:        payload.ImageURL,
	})
	if err != nil {
		return nil, err
	}

	o.sendConfirmationEmail(ctx, sig)

	return sig, nil
}

func (o *Operations) CreateRequest(ctx context.Context, payload models.CreateSignatureRequestInput) (*models.SignatureRequest, error) {
	edition, err := o.Editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	ident, err := o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	idempotencyKey := uuid.New()
	expiresAt := time.Now().Add(time.Duration(payload.ExpiresInDays) * 24 * time.Hour)
	request, err := o.Requests.CreateRequest(ctx, &models.SignatureRequest{
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

	event, err := o.Events.GetByID(ctx, edition.EventID)
	if err != nil {
		return nil, err
	}

	token, err := crypto.SignHMACJWT(models.SignatureRequestClaims{
		RequestID: request.ID,
		EditionID: payload.EditionID,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}, []byte(o.HmacSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to encode signature token: %w", err)
	}

	link := fmt.Sprintf("%s/signature-requests/fulfill?token=%s", o.AppURL, token)
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

	err = o.Email.Send(ctx, email.Message{
		To:      []string{*payload.SignatoryEmail},
		Subject: fmt.Sprintf("Signature request for %s — %s", event.FullName, edition.Name),
		Body:    body,
		HTML:    true,
	})
	if err != nil {
		telemetry.Log().Error("failed to send signature request email", telemetry.Err(err))
	}

	return request, nil
}

func (o *Operations) Delete(ctx context.Context, id uuid.UUID) error {
	existing, err := o.Signatures.GetByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, existing.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Signatures.Delete(ctx, id)
}

func (o *Operations) DenyRequest(ctx context.Context, token string, reason *string) error {
	var claims models.SignatureRequestClaims
	_, err := crypto.ParseHMACJWT(token, &claims, []byte(o.HmacSecret))
	if err != nil {
		return errs.ErrInvalidSignatureToken
	}

	request, err := o.Requests.GetRequestByID(ctx, claims.RequestID)
	if err != nil {
		return err
	}
	if request.Status != models.SignatureRequestStatusPending {
		return errs.ErrSignatureRequestNotPending
	}
	if request.EditionID != claims.EditionID {
		return errs.ErrTokenRequestEditionMismatch
	}

	return o.Requests.CancelRequest(ctx, request.ID, reason)
}

func (o *Operations) FulfillRequest(ctx context.Context, token string, imageURL string) (*models.Signature, error) {
	var claims models.SignatureRequestClaims
	_, err := crypto.ParseHMACJWT(token, &claims, []byte(o.HmacSecret))
	if err != nil {
		return nil, errs.ErrInvalidSignatureToken
	}

	request, err := o.Requests.GetRequestByID(ctx, claims.RequestID)
	if err != nil {
		return nil, err
	}
	if request.Status != models.SignatureRequestStatusPending {
		return nil, errs.ErrSignatureRequestNotPending
	}
	if request.EditionID != claims.EditionID {
		return nil, errs.ErrTokenRequestEditionMismatch
	}

	sig, err := o.Signatures.Create(ctx, &models.Signature{
		EditionID:       request.EditionID,
		CreatedBy:       request.CreatedBy,
		SignatoryName:   request.SignatoryName,
		SignatoryTitle:  request.SignatoryTitle,
		SignatoryEmail:  request.SignatoryEmail,
		SignatoryUserID: request.SignatoryUserID,
		ImageURL:        imageURL,
	})
	if err != nil {
		return nil, err
	}

	err = o.Requests.CompleteRequest(ctx, request.ID, sig.ID)
	if err != nil {
		return nil, err
	}

	o.sendConfirmationEmail(ctx, sig)

	return sig, nil
}

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.Signature, error) {
	return o.Signatures.GetByID(ctx, id)
}

func (o *Operations) GetRequestByID(ctx context.Context, id uuid.UUID) (*models.SignatureRequest, error) {
	return o.Requests.GetRequestByID(ctx, id)
}

func (o *Operations) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Signature, error) {
	return o.Signatures.ListByEdition(ctx, editionID)
}

func (o *Operations) ListRequestsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.SignatureRequest, error) {
	return o.Requests.ListRequestsByEdition(ctx, editionID)
}

func (o *Operations) RevokeSignature(ctx context.Context, token string) error {
	var claims models.SignatureRevocationClaims
	_, err := crypto.ParseHMACJWT(token, &claims, []byte(o.HmacSecret))
	if err != nil {
		return errs.ErrInvalidSignatureToken
	}

	sig, err := o.Signatures.GetByID(ctx, claims.SignatureID)
	if err != nil {
		return err
	}
	if sig.EditionID != claims.EditionID {
		return errs.ErrTokenSignatureEditionMismatch
	}

	return o.Signatures.Delete(ctx, sig.ID)
}

func (o *Operations) sendConfirmationEmail(ctx context.Context, sig *models.Signature) {
	if sig.SignatoryEmail == nil {
		return
	}

	edition, err := o.Editions.GetByID(ctx, sig.EditionID)
	if err != nil {
		return
	}

	event, err := o.Events.GetByID(ctx, edition.EventID)
	if err != nil {
		return
	}

	token, err := crypto.SignHMACJWT(models.SignatureRevocationClaims{
		SignatureID: sig.ID,
		EditionID:   sig.EditionID,
		IssuedAt:    jwt.NewNumericDate(time.Now()),
	}, []byte(o.HmacSecret))
	if err != nil {
		return
	}

	revokeLink := fmt.Sprintf("%s/signatures/revoke?token=%s", o.AppURL, token)
	body, err := assets.RenderSignatureCreatedEmail(assets.SignatureCreatedEmailData{
		SignatoryName: sig.SignatoryName,
		EventName:     event.FullName,
		EditionName:   edition.Name,
		RevokeLink:    revokeLink,
	})
	if err != nil {
		return
	}

	err = o.Email.Send(ctx, email.Message{
		To:      []string{*sig.SignatoryEmail},
		Subject: fmt.Sprintf("Your signature for %s — %s", event.FullName, edition.Name),
		Body:    body,
		HTML:    true,
	})
	if err != nil {
		telemetry.Log().Error("failed to send signature confirmation email", telemetry.Err(err))
	}
}
