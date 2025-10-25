package service

import (
	"time"

	"GoAuth/internal/models"
	"GoAuth/internal/repository"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/golang-jwt/jwt/v5"
  "github.com/google/uuid"
	"github.com/spf13/viper"
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
			ID: accessJTI,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessJTIID, _ := uuid.Parse(accessJTI)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := accessToken.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", uuid.Nil, resp.InternalServerError("error signing access token").WithTracePrefix("error").AddTrace(err)
	}
	return tokenStr, accessJTIID, nil
}

func newRefreshToken(accessJTI, refreshJTI uuid.UUID, agent, ip string, expires_at time.Time, session_id uuid.UUID) (string, *resp.Response) {
	claims := models.RefreshClaims{
		Sub: models.RefreshSubJWT{
			AccessJTI: accessJTI,
			UserAgent: agent,
			UserIP: ip,
      SessionID: session_id,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires_at),
			Issuer:    "GoAuth",
			ID: refreshJTI.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := refreshToken.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", resp.InternalServerError("error signing refresh token").WithTracePrefix("error").AddTrace(err)
	}
	return tokenStr, nil
}
