package auth

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	authport "GoAuth/internal/ports/auth"
	"GoAuth/internal/ports/outbound"
	"GoAuth/internal/utils"
	"context"
	"crypto/ed25519"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
			return nil, apierr.ErrUnauthorized.WithMsg("missing kid").WithID(apierr.TokenMissingKid)
		}

		return uc.resolvePublicKey(ctx, kid)
	})

	if err != nil {
		return nil, apierr.FromJWTError(err, "access token")
	}

	if !token.Valid {
		return nil, apierr.ErrUnauthorized.WithMsg("invalid access token").WithID(apierr.TokenInvalid)
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
			return nil, apierr.ErrUnauthorized.WithMsg("missing kid").WithID(apierr.TokenMissingKid)
		}

		return uc.resolvePublicKey(ctx, kid)
	})

	if err != nil {
		return nil, apierr.FromJWTError(err, "refresh token")
	}

	if !token.Valid {
		return nil, apierr.ErrUnauthorized.WithMsg("invalid refresh token").WithID(apierr.TokenInvalid)
	}

	return claims, nil
}

func (uc *TokenVerifier) resolvePublicKey(ctx context.Context, kid string) (ed25519.PublicKey, error) {
	switch {
	case kid == "goauth:v1":
		return utils.GoAuthPublicKey, nil

	case strings.HasPrefix(kid, "project:"):
		parts := strings.Split(kid, ":")
		if len(parts) != 3 {
			return nil, apierr.ErrInvalidInput.WithMsg("invalid token kid").WithID(apierr.TokenInvalidKid)
		}

		projectID, err := uuid.Parse(parts[1])
		if err != nil {
			return nil, apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
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
			return nil, apierr.ErrInternal.WithMsg("failed to parse project public key").WithCause(err).WithID(apierr.ProjectFailedToParseKey)
		}

		return decodedKey, nil

	default:
		return nil, apierr.ErrUnauthorized.WithMsg("unknown kid").WithID(apierr.TokenUnknownKid)
	}
}
