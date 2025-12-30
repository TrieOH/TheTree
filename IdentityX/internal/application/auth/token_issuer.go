package auth

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/user"
	"GoAuth/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func newAccessToken(user user.User, ip, agent string, sessionId uuid.UUID, expiresAt time.Time) (string, uuid.UUID, error) {
	accessJTI := uuid.NewString()
	claims := auth.AccessClaims{
		Sub: auth.AccessSubJWT{
			ID:        user.ID,
			UserType:  user.UserType,
			Email:     user.Email,
			SessionID: sessionId,
			UserAgent: agent,
			UserIP:    ip,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "GoAuth",
			ID:        accessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessJTIID, _ := uuid.Parse(accessJTI)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenStr, err := accessToken.SignedString(utils.GoAuthPrivateKey)
	if err != nil {
		return "", uuid.Nil, apierr.ErrInternal.WithMsg("error signing access token").WithID(apierr.TokenCouldNotSign).WithCause(err)
	}
	return tokenStr, accessJTIID, nil
}

func newRefreshToken(accessJTI, refreshJTI uuid.UUID, expiresAt time.Time) (string, error) {
	claims := auth.RefreshClaims{
		Sub: auth.RefreshSubJWT{
			AccessJTI: accessJTI,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "GoAuth",
			ID:        refreshJTI.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenStr, err := refreshToken.SignedString(utils.GoAuthPrivateKey)
	if err != nil {
		return "", apierr.ErrInternal.WithMsg("error signing refresh token").WithID(apierr.TokenCouldNotSign).WithCause(err)
	}
	return tokenStr, nil
}

func newProjectAccessToken(user project_users.ProjectUser, ip, agent string, sessionId uuid.UUID, expiresAt time.Time) (string, uuid.UUID, error) {
	accessJTI := uuid.NewString()
	claims := auth.AccessClaims{
		Sub: auth.AccessSubJWT{
			ID:        user.ID,
			UserType:  user.UserType,
			ProjectID: &user.ProjectID,
			Metadata:  user.Metadata,
			Email:     user.Email,
			SessionID: sessionId,
			UserAgent: agent,
			UserIP:    ip,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "GoAuth",
			ID:        accessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessJTIID, _ := uuid.Parse(accessJTI)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenStr, err := accessToken.SignedString(utils.GoAuthPrivateKey)
	if err != nil {
		return "", uuid.Nil, apierr.ErrInternal.WithMsg("error signing access token").WithID(apierr.TokenCouldNotSign).WithCause(err)
	}
	return tokenStr, accessJTIID, nil
}
