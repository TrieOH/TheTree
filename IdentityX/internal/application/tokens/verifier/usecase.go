package verifier

import (
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/ports/inbounds"
	"context"
	"encoding/base64"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type TokenVerifier struct {
	keys inbounds.KeysService
}

func NewTokenVerifier(keys inbounds.KeysService) inbounds.TokenVerifier {
	return &TokenVerifier{keys: keys}
}

var _ inbounds.TokenVerifier = (*TokenVerifier)(nil)

func (uc *TokenVerifier) VerifyAccessToken(
	ctx context.Context,
	tokenStr string,
) (*auth.AccessClaims, error) {
	return verifyToken(
		ctx,
		uc,
		"access",
		tokenStr,
		&auth.AccessClaims{},
	)
}

func (uc *TokenVerifier) VerifyRefreshToken(
	ctx context.Context,
	tokenStr string,
) (*auth.RefreshClaims, error) {
	return verifyToken(
		ctx,
		uc,
		"refresh",
		tokenStr,
		&auth.RefreshClaims{},
	)
}

func (uc *TokenVerifier) VerifyVerificationToken(
	ctx context.Context,
	tokenStr string,
) (*auth.VerificationClaims, error) {
	return verifyToken(
		ctx,
		uc,
		"verification",
		tokenStr,
		&auth.VerificationClaims{},
	)
}

func verifyToken[T jwt.Claims](
	ctx context.Context,
	uc *TokenVerifier,
	tokenType string,
	tokenStr string,
	claims T,
) (T, error) {
	token, err := parseJWTUnverified(tokenStr, claims)
	if err != nil {
		return claims, apierr.FromJWTError(err, tokenType)
	}

	alg, _ := token.Header["alg"].(string)
	if alg != jwt.SigningMethodEdDSA.Alg() {
		return claims, auth.ErrTokenInvalidAlg{
			TokenType: tokenType,
			Expected:  jwt.SigningMethodEdDSA.Alg(),
			Got:       alg,
		}
	}

	if token.Method == nil || token.Method.Alg() != jwt.SigningMethodEdDSA.Alg() {
		methodAlg := ""
		if token.Method != nil {
			methodAlg = token.Method.Alg()
		}
		return claims, auth.ErrTokenInvalidAlg{
			TokenType: tokenType,
			Expected:  jwt.SigningMethodEdDSA.Alg(),
			Got:       methodAlg,
		}
	}

	kid, ok := token.Header["kid"].(string)
	if !ok || kid == "" {
		return claims, auth.ErrTokenMissingKID{TokenType: tokenType}
	}

	payload, sig, err := splitJWT(tokenType, tokenStr)
	if err != nil {
		return claims, err
	}

	switch {
	case strings.HasPrefix(kid, "goauth:"):
		parts := strings.Split(kid, ":")
		if len(parts) < 2 {
			return claims, auth.ErrTokenInvalidKID{TokenType: tokenType}
		}

		if err := uc.keys.VerifyGoAuth(ctx, kid, payload, sig); err != nil {
			if apierr.IsNotFound(err) {
				return claims, auth.ErrTokenUntrusted{TokenType: tokenType}
			}
			return claims, err
		}

	case strings.HasPrefix(kid, "project:"):
		parts := strings.Split(kid, ":")
		if len(parts) < 3 {
			return claims, auth.ErrTokenInvalidKID{TokenType: tokenType}
		}

		projectID, err := validation.ParseUUID(parts[1], "project_id")
		if err != nil {
			return claims, err
		}

		if err := uc.keys.VerifyProject(ctx, projectID, kid, payload, sig); err != nil {
			if apierr.IsNotFound(err) {
				return claims, auth.ErrTokenUntrusted{TokenType: tokenType}
			}
			return claims, err
		}

	default:
		return claims, auth.ErrTokenUnknownKID{TokenType: tokenType}
	}

	return claims, nil
}

func parseJWTUnverified[T jwt.Claims](tokenStr string, claims T) (*jwt.Token, error) {
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(tokenStr, claims)
	return token, err
}

func splitJWT(tokenType, tokenStr string) (signingInput, sig []byte, err error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, nil, auth.ErrTokenInvalidFormat{TokenType: tokenType}
	}

	signingInput = []byte(parts[0] + "." + parts[1])

	sig, err = base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, nil, err
	}

	return signingInput, sig, nil
}
