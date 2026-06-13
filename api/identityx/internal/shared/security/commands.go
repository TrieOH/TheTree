package security

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"IdentityX/contracts"
	"lib/crypto"
	"lib/errx"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func SignKey(payload []byte, pair *contracts.Pair, encryptionKey []byte) ([]byte, error) {
	decrypted, err := crypto.Decrypt(pair.PrivateKey, encryptionKey)
	if err != nil {
		return nil, err
	}

	priv, err := parseEd25519Private(decrypted)
	if err != nil {
		return nil, err
	}

	return ed25519.Sign(priv, payload), nil
}

func VerifyKeyPair(projectID *uuid.UUID, payload, sig []byte, pair *contracts.Pair) error {
	if projectID != nil {
		if pair.ProjectID == nil || *pair.ProjectID != *projectID {
			return fun.ErrUnauthorized("project keys mismatch")
		}
	}

	pub, err := parseEd25519Public(pair.PublicKey)
	if err != nil {
		return err
	}

	if !ed25519.Verify(pub, payload, sig) {
		return fun.ErrUnauthorized("invalid key signature")
	}

	return nil
}

func parseEd25519Private(pemBytes []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid PEM Private Key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing PKCS8 Public key: %w", err)
	}

	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not an ED25519 private key")
	}

	return priv, nil
}

func parseEd25519Public(pemStr string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("invalid PEM Public Key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing PKIX Public key: %w", err)
	}

	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("public key is not an ED25519 public key")
	}

	return pub, nil
}

func NewAccessToken(in contracts.NewAccessTokenInput, issuer string) ([]byte, error) {
	if in.User.ProjectID != nil {
		issuer = in.User.ProjectID.String()
	}

	claims := contracts.AccessClaims{
		Sub: contracts.AccessSub{
			ID:         in.User.ID,
			UserType:   string(in.User.UserType),
			ProjectID:  in.User.ProjectID,
			Email:      in.User.Email,
			SessionID:  in.SessionID,
			UserAgent:  in.Agent,
			UserIP:     in.IP,
			IsVerified: in.User.IsVerified,
			FamilyID:   in.FamilyID,
			VerifiedAt: in.User.VerifiedAt,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    issuer,
			ID:        in.AccessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = in.KID

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func NewRefreshToken(in contracts.NewRefreshTokenInput, issuer string) ([]byte, error) {
	claims := contracts.RefreshClaims{
		Sub: contracts.RefreshSub{
			AccessJTI: in.AccessJTI,
			FamilyID:  in.FamilyID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    issuer,
			ID:        in.RefreshJTI.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = in.KID

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func NewVerificationToken(in contracts.NewVerificationTokenInput, issuer string) ([]byte, error) {
	now := time.Now()
	claims := contracts.VerificationClaims{
		Sub: contracts.VerificationSub{
			Subject: in.Subject,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-1 * time.Minute)),
			Audience:  jwt.ClaimStrings{"email-verification"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = in.KID

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func NewResetPasswordToken(in contracts.NewResetPasswordInput, issuer string) ([]byte, error) {
	now := time.Now()
	claims := contracts.ResetPasswordClaims{
		Sub: contracts.ResetPasswordSub{
			Subject:   in.Subject,
			ProjectID: in.ProjectID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(in.ExpiresAt),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-1 * time.Minute)),
			Audience:  jwt.ClaimStrings{"password-reset"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = in.KID

	payload, err := token.SigningString()
	if err != nil {
		return nil, err
	}

	return []byte(payload), nil
}

func AssembleJWT(payload []byte, sig []byte) string {
	return string(payload) + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func VerifyAccessToken(
	tokenStr string,
	pair *contracts.Pair,
) (*contracts.AccessClaims, error) {
	return verifyToken(
		pair,
		"access",
		tokenStr,
		&contracts.AccessClaims{},
	)
}

func VerifyRefreshToken(
	tokenStr string,
	pair *contracts.Pair,
) (*contracts.RefreshClaims, error) {
	return verifyToken(
		pair,
		"refresh",
		tokenStr,
		&contracts.RefreshClaims{},
	)
}

func VerifyVerificationToken(
	tokenStr string,
	pair *contracts.Pair,
) (*contracts.VerificationClaims, error) {
	return verifyToken(
		pair,
		"verification",
		tokenStr,
		&contracts.VerificationClaims{},
	)
}

func VerifyResetPasswordToken(
	tokenStr string,
	pair *contracts.Pair,
) (*contracts.ResetPasswordClaims, error) {
	return verifyToken(
		pair,
		"reset password",
		tokenStr,
		&contracts.ResetPasswordClaims{},
	)
}

func verifyToken[T jwt.Claims](
	pair *contracts.Pair,
	tokenType string,
	tokenStr string,
	claims T,
) (T, error) {
	token, err := ParseJWTUnverified(tokenStr, claims)
	if err != nil {
		return claims, errx.FromJWTError(err, tokenType)
	}

	alg, _ := token.Header["alg"].(string)
	if alg != jwt.SigningMethodEdDSA.Alg() {
		return claims, fun.Errf("invalid %s token alg, expected (%s) but got (%s)", tokenType, jwt.SigningMethodEdDSA.Alg(), alg).Unauthorized()
	}

	if token.Method == nil || token.Method.Alg() != jwt.SigningMethodEdDSA.Alg() {
		methodAlg := ""
		if token.Method != nil {
			methodAlg = token.Method.Alg()
		}
		return claims, fun.Errf("invalid %s token alg, expected (%s) but got (%s)", tokenType, jwt.SigningMethodEdDSA.Alg(), methodAlg).Unauthorized()
	}

	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return claims, fun.Errf("%s token missing kid", tokenType).Unauthorized()
	}
	payload, sig, err := splitJWT(tokenType, tokenStr)
	if err != nil {
		return claims, err
	}

	switch {
	case strings.HasPrefix(kid, "goauth:"):
		parts := strings.Split(kid, ":")
		if len(parts) < 2 {
			return claims, fun.Errf("invalid %s token kid", tokenType).Unauthorized()
		}
		if err = VerifyKeyPair(nil, payload, sig, pair); err != nil {
			return claims, err
		}

	case strings.HasPrefix(kid, "project:"):
		parts := strings.Split(kid, ":")
		if len(parts) < 3 {
			return claims, fun.Errf("invalid %s token kid", tokenType).Unauthorized()
		}
		var projectID uuid.UUID
		projectID, err = uuid.Parse(parts[1])
		if err != nil {
			return claims, fun.Errf("invalid project kid project: %s", err).Unauthorized()
		}
		if err = VerifyKeyPair(&projectID, payload, sig, pair); err != nil {
			return claims, err
		}

	default:
		return claims, fun.Errf("unknown %s token kid", tokenType).Unauthorized()
	}

	return claims, nil
}

func ParseJWTUnverified[T jwt.Claims](tokenStr string, claims T) (*jwt.Token, error) {
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(tokenStr, claims)
	return token, err
}

func splitJWT(tokenType, tokenStr string) (signingInput, sig []byte, err error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, nil, fun.Errf("invalid %s token format", tokenType).Unauthorized()
	}
	signingInput = []byte(parts[0] + "." + parts[1])
	sig, err = base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, nil, err
	}
	return signingInput, sig, nil
}
