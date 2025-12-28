package service

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/utils"
	"time"

	"GoAuth/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func newAccessToken(user models.User, ip, agent string, sessionId uuid.UUID) (string, uuid.UUID, error) {
	accessJTI := uuid.NewString()
	claims := models.AccessClaims{
		Sub: models.AccessSubJWT{
			ID:        user.ID,
			UserType:  user.UserType,
			Email:     user.Email,
			SessionID: sessionId,
			UserAgent: agent,
			UserIP:    ip,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
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
	claims := models.RefreshClaims{
		Sub: models.RefreshSubJWT{
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

func newProjectAccessToken(user models.ProjectUser, ip, agent string, sessionId uuid.UUID) (string, uuid.UUID, error) {
	accessJTI := uuid.NewString()
	claims := models.AccessClaims{
		Sub: models.AccessSubJWT{
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
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
