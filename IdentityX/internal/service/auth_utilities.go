package service

import (
	"time"

	"GoAuth/internal/models"
	"GoAuth/internal/repository"

	resp "github.com/MintzyG/GoResponse/response"
	"github.com/golang-jwt/jwt/v5"
  "github.com/google/uuid"
	"github.com/spf13/viper"
)

func newAccessToken(dbUser repository.User) (string, *resp.Response) {
	claims := models.AccessClaims{
		Sub: models.AccessSubJWT{
			ID:    dbUser.ID,
			Email: dbUser.Email,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Issuer:    "GoAuth",
			ID: uuid.NewString(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := accessToken.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", resp.InternalServerError("error signing access token").WithTracePrefix("error").AddTrace(err)
	}
	return tokenStr, nil
}

func newRefreshToken() (string, *resp.Response) {
	claims := models.RefreshClaims{
		Sub: models.RefreshSubJWT{
			MetaData: "Refresh Metadata",
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			Issuer:    "GoAuth",
			ID: uuid.NewString(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := refreshToken.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", resp.InternalServerError("error signing refresh token").WithTracePrefix("error").AddTrace(err)
	}
	return tokenStr, nil
}
