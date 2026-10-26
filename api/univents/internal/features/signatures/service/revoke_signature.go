package service

import (
	"context"
	"lib/crypto"
	"univents/errs"
	"univents/models"
)

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
