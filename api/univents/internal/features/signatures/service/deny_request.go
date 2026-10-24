package service

import (
	"context"
	"lib/crypto"
	"lib/telemetry"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (o *Operations) DenyRequest(ctx context.Context, token string, reason *string) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.DenyRequest")
	defer span.End()

	var claims models.SignatureRequestClaims
	_, err := crypto.ParseHMACJWT(token, &claims, []byte(o.HmacSecret))
	if err != nil {
		return fun.ErrForbidden("invalid or expired token")
	}

	request, err := o.Requests.GetRequestByID(ctx, claims.RequestID)
	if err != nil {
		return err
	}
	if request.Status != models.SignatureRequestStatusPending {
		return fun.ErrConflict("signature request is no longer pending")
	}
	if request.EditionID != claims.EditionID {
		return fun.ErrForbidden("token does not match request edition")
	}

	return o.Requests.CancelRequest(ctx, request.ID, reason)
}
