package service

import (
	"GoAuth/internal/utils"
	"time"

	"GoAuth/internal/models"
	"GoAuth/internal/repository"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func newAccessToken(dbUser repository.User) (string, uuid.UUID, *resp.Response) {
	accessJTI := uuid.NewString()
	claims := models.AccessClaims{
		Sub: models.AccessSubJWT{
			ID:    dbUser.ID,
			Email: dbUser.Email,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    "GoAuth",
			ID:        accessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessJTIID, _ := uuid.Parse(accessJTI)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenStr, err := accessToken.SignedString(utils.GoAuthPrivateKey)
	if err != nil {
		return "", uuid.Nil, resp.InternalServerError("error signing access token").WithTracePrefix("error").AddTrace(err)
	}
	return tokenStr, accessJTIID, nil
}

func newRefreshToken(accessJTI, refreshJTI uuid.UUID, agent, ip string, expiresAt time.Time, sessionId uuid.UUID) (string, *resp.Response) {
	claims := models.RefreshClaims{
		Sub: models.RefreshSubJWT{
			AccessJTI: accessJTI,
			UserAgent: agent,
			UserIP:    ip,
			SessionID: sessionId,
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
		return "", resp.InternalServerError("error signing refresh token").WithTracePrefix("error").AddTrace(err)
	}
	return tokenStr, nil
}
