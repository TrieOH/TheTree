package issuer

import (
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/user"
	"GoAuth/internal/ports/inbounds"
	"crypto/ed25519"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type UseCase struct{}

var _ inbounds.TokenIssuer = (*UseCase)(nil)

func NewTokenIssuer() inbounds.TokenIssuer {
	return &UseCase{}
}

func (uc *UseCase) NewAccessToken(user user.User, key ed25519.PrivateKey, ip, agent, accessJTI, keyID string, sessionID uuid.UUID, expiresAt time.Time) (string, error) {
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
		return "", auth.ErrSigningToken{TokenType: "access", Cause: err}
	}
	return tokenStr, nil
}

func (uc *UseCase) NewRefreshToken(keyID string, privKey ed25519.PrivateKey, accessJTI, refreshJTI uuid.UUID, expiresAt time.Time) (string, error) {
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
		return "", auth.ErrSigningToken{TokenType: "refresh", Cause: err}
	}
	return tokenStr, nil
}

func (uc *UseCase) NewProjectAccessToken(user project_users.ProjectUser, ip, agent, accessJTI, keyID string, sessionID uuid.UUID, expiresAt time.Time, privKey ed25519.PrivateKey) (string, error) {
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
		return "", auth.ErrSigningToken{TokenType: "access", Cause: err}
	}
	return tokenStr, nil
}
