package commands

import (
	"context"
	"lib/crypto"
	"lib/telemetry"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (c *Commands) RevokeSignature(ctx context.Context, token string) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.RevokeSignature")
	defer span.End()

	var claims models.SignatureRevocationClaims
	_, err := crypto.ParseHMACJWT(token, &claims, []byte(c.hmacSecret))
	if err != nil {
		return fun.ErrForbidden("invalid or expired token")
	}

	sig, err := c.signatures.GetByID(ctx, claims.SignatureID)
	if err != nil {
		return err
	}
	if sig.EditionID != claims.EditionID {
		return fun.ErrForbidden("token does not match signature edition")
	}

	return c.signatures.Delete(ctx, sig.ID)
}
