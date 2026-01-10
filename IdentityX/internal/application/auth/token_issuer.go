package auth

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/user"
	"crypto/ed25519"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func newAccessToken(user user.User, key ed25519.PrivateKey, ip, agent, accessJTI, keyID string, sessionID uuid.UUID, expiresAt time.Time) (string, error) {
	claims := auth.AccessClaims{
		Sub: auth.AccessSubJWT{
			ID:        user.ID,
			UserType:  user.UserType,
			Email:     user.Email,
			SessionID: sessionID,
			UserAgent: agent,
			UserIP:    ip,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    viper.GetString("ISSUER"),
			ID:        accessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	accessToken.Header["kid"] = keyID
	tokenStr, err := accessToken.SignedString(key)
	if err != nil {
		return "", apierr.ErrInternal.WithMsg("error signing access token").WithID(apierr.TokenCouldNotSign).WithCause(err)
	}
	return tokenStr, nil
}

func newRefreshToken(keyID string, privKey ed25519.PrivateKey, accessJTI, refreshJTI uuid.UUID, expiresAt time.Time) (string, error) {
	claims := auth.RefreshClaims{
		Sub: auth.RefreshSubJWT{
			AccessJTI: accessJTI,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    viper.GetString("ISSUER"),
			ID:        refreshJTI.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	refreshToken.Header["kid"] = keyID
	tokenStr, err := refreshToken.SignedString(privKey)
	if err != nil {
		return "", apierr.ErrInternal.WithMsg("error signing refresh token").WithID(apierr.TokenCouldNotSign).WithCause(err)
	}
	return tokenStr, nil
}

func newProjectAccessToken(user project_users.ProjectUser, ip, agent, accessJTI, keyID string, sessionID uuid.UUID, expiresAt time.Time, privKey ed25519.PrivateKey) (string, error) {
	claims := auth.AccessClaims{
		Sub: auth.AccessSubJWT{
			ID:        user.ID,
			UserType:  user.UserType,
			ProjectID: &user.ProjectID,
			Metadata:  user.Metadata,
			Email:     user.Email,
			SessionID: sessionID,
			UserAgent: agent,
			UserIP:    ip,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    viper.GetString("ISSUER"),
			ID:        accessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	accessToken.Header["kid"] = keyID
	tokenStr, err := accessToken.SignedString(privKey)
	if err != nil {
		return "", apierr.ErrInternal.WithMsg("error signing access token").WithID(apierr.TokenCouldNotSign).WithCause(err)
	}
	return tokenStr, nil
}
