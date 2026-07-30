package commands

import (
	"context"
	"lib/crypto"
	"lib/telemetry"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) FulfillRequest(ctx context.Context, token string, imageURL string) (*models.Signature, error) {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.FulfillRequest")
	defer span.End()

	var claims models.SignatureRequestClaims
	_, err := crypto.ParseHMACJWT(token, &claims, []byte(c.hmacSecret))
	if err != nil {
		return nil, fun.ErrForbidden("invalid or expired token")
	}

	request, err := c.requests.GetRequestByID(ctx, claims.RequestID)
	if err != nil {
		return nil, err
	}
	if request.Status != models.SignatureRequestStatusPending {
		return nil, fun.ErrConflict("signature request is no longer pending")
	}
	if request.EditionID != claims.EditionID {
		return nil, fun.ErrForbidden("token does not match request edition")
	}

	sig, err := c.signatures.Create(ctx, &models.Signature{
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

	err = c.requests.CompleteRequest(ctx, request.ID, sig.ID)
	if err != nil {
		return nil, err
	}

	c.sendConfirmationEmail(ctx, sig)

	return sig, nil
}
