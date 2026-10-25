package service

import (
	"context"
	"lib/crypto"
	"univents/errs"
	"univents/models"
)

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
