package service

import (
	"context"
	"lib/crypto"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (o *Operations) RevokeSignature(ctx context.Context, token string) error {
	var claims models.SignatureRevocationClaims
	_, err := crypto.ParseHMACJWT(token, &claims, []byte(o.HmacSecret))
	if err != nil {
		return fun.ErrForbidden("invalid or expired token")
	}

	sig, err := o.Signatures.GetByID(ctx, claims.SignatureID)
	if err != nil {
		return err
	}
	if sig.EditionID != claims.EditionID {
		return fun.ErrForbidden("token does not match signature edition")
	}

	return o.Signatures.Delete(ctx, sig.ID)
}
