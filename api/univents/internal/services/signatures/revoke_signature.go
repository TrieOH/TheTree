package signatures

import (
	"context"
	"lib/crypto"
	"lib/telemetry"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (o *Operations) RevokeSignature(ctx context.Context, token string) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.RevokeSignature")
	defer span.End()

	var claims models.SignatureRevocationClaims
	_, err := crypto.ParseHMACJWT(token, &claims, []byte(o.hmacSecret))
	if err != nil {
		return fun.ErrForbidden("invalid or expired token")
	}

	sig, err := o.signatures.GetByID(ctx, claims.SignatureID)
	if err != nil {
		return err
	}
	if sig.EditionID != claims.EditionID {
		return fun.ErrForbidden("token does not match signature edition")
	}

	return o.signatures.Delete(ctx, sig.ID)
}
