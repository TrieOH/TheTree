package auth

import (
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	authport "GoAuth/internal/ports/auth"
	"GoAuth/internal/ports/outbound"
	"GoAuth/internal/utils"
	"context"
	"crypto/ed25519"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type TokenVerifier struct {
	projects outbound.ProjectRepository
}

func NewTokenVerifier(projects outbound.ProjectRepository) authport.TokenVerifier {
	return &TokenVerifier{projects: projects}
}

var _ authport.TokenVerifier = (*TokenVerifier)(nil)

func (uc *TokenVerifier) VerifyAccessToken(
	ctx context.Context,
	tokenStr string,
) (*auth.AccessClaims, error) {

	claims := &auth.AccessClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, auth.ErrTokenMissingKID{TokenType: "access"}
		}

		return uc.resolvePublicKey(ctx, kid, "access")
	})

	if err != nil {
		return nil, apierr.FromJWTError(err, "access")
	}

	if !token.Valid {
		return nil, auth.ErrInvalidToken{TokenType: "access"}
	}

	return claims, nil
}

func (uc *TokenVerifier) VerifyRefreshToken(
	ctx context.Context,
	tokenStr string,
) (*auth.RefreshClaims, error) {

	claims := &auth.RefreshClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, auth.ErrTokenMissingKID{TokenType: "refresh"}
		}

		return uc.resolvePublicKey(ctx, kid, "refresh")
	})

	if err != nil {
		return nil, apierr.FromJWTError(err, "refresh")
	}

	if !token.Valid {
		return nil, auth.ErrInvalidToken{TokenType: "refresh"}
	}

	return claims, nil
}

func (uc *TokenVerifier) resolvePublicKey(ctx context.Context, kid, tokenType string) (ed25519.PublicKey, error) {
	switch {
	case kid == "goauth:v1":
		return utils.GoAuthPublicKey, nil

	case strings.HasPrefix(kid, "project:"):
		parts := strings.Split(kid, ":")
		if len(parts) != 3 {
			return nil, auth.ErrTokenInvalidKID{TokenType: tokenType}
		}

		projectID, err := validation.ParseUUID(parts[1], "project_id")
		if err != nil {
			return nil, err
		}

		// TODO: Implement key rotation for projects, only then start using versioned keys
		// keyVersion = parts[2]

		pubKey, err := uc.projects.GetPublicKeyByID(ctx, projectID)
		if err != nil {
			return nil, err
		}
		var decodedKey ed25519.PublicKey
		decodedKey, err = utils.ParseEd25519PublicKey(pubKey)
		if err != nil {
			return nil, utils.ErrParseProjectKey{KeyType: "public", Cause: err}
		}
		return decodedKey, nil
	default:
		return nil, auth.ErrTokenUnknownKID{TokenType: tokenType}
	}
}
