// Package tokens owns the shared token-verification reasoning for
// IdentityX: opening a signed token, resolving its signing key by kid
// (rejecting revoked keys), verifying the signature, and rejecting
// blacklisted tokens. The auth middleware and the authn rotation
// operations (refresh, logout) cross the same seam, so key and revocation
// policy live in one module instead of being re-derived per caller.
package tokens

import (
	"context"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Verifier resolves signing keys and verifies tokens against them.
type Verifier struct {
	cryptoKeys ports.CryptoKeysRepo
	blacklist  ports.BlacklistRepo
}

func NewVerifier(cryptoKeys ports.CryptoKeysRepo, blacklist ports.BlacklistRepo) *Verifier {
	return &Verifier{cryptoKeys: cryptoKeys, blacklist: blacklist}
}

// Verify opens tokenStr, resolves the signing key by kid (revoked keys are
// rejected), verifies the signature, and rejects blacklisted tokens. The
// parsed claims are populated into claims; the signing key is returned for
// callers that need it.
func (v *Verifier) Verify(ctx context.Context, tokenStr string, claims jwt.Claims) (*jwt.Token, *models.CryptoKey, error) {
	token, err := crypto.OpenUnverified(tokenStr, claims)
	if err != nil {
		return nil, nil, err
	}

	key, err := v.KeyForToken(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	_, err = crypto.VerifyToken(tokenStr, key.PublicKey, claims)
	if err != nil {
		return nil, nil, fun.ErrUnauthorized("invalid token")
	}

	err = v.ensureNotBlacklisted(ctx, tokenJTI(claims))
	if err != nil {
		return nil, nil, err
	}

	return token, key, nil
}

// tokenJTI extracts the jti (RegisteredClaims.ID) from the two claim types
// IdentityX signs; the jwt.Claims interface does not expose it.
func tokenJTI(claims jwt.Claims) string {
	switch c := claims.(type) {
	case *models.AccessClaims:
		return c.ID
	case *models.RefreshClaims:
		return c.ID
	default:
		return ""
	}
}

// KeyForToken resolves the signing key for an already-opened token,
// rejecting missing and revoked keys. For flows that open a token without
// verifying it (logout's access token) but still need its signing key.
func (v *Verifier) KeyForToken(ctx context.Context, token *jwt.Token) (*models.CryptoKey, error) {
	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return nil, fun.ErrUnauthorized("missing kid")
	}
	keyID, err := uuid.Parse(kid)
	if err != nil {
		return nil, fun.ErrUnauthorized("invalid kid")
	}
	key, err := v.cryptoKeys.GetByID(ctx, keyID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrUnauthorized("outdated token")
	}
	if err != nil {
		return nil, err
	}
	if key.Status == "revoked" {
		return nil, fun.ErrUnauthorized("token signing key revoked")
	}
	return key, nil
}

// ensureNotBlacklisted rejects tokens revoked at logout or refresh
// rotation. A missing entry (never revoked) passes.
func (v *Verifier) ensureNotBlacklisted(ctx context.Context, jti string) error {
	_, err := v.blacklist.GetByTargetAndType(ctx, jti, models.BlacklistEntryTypeToken)
	if err == nil {
		return fun.ErrUnauthorized("token has been revoked")
	}
	if !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	return nil
}
