package service

import (
	"time"

	"GoAuth/internal/models"
	"GoAuth/internal/repository"

	resp "github.com/MintzyG/GoResponse/response"
	"github.com/golang-jwt/jwt/v5"
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
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := refreshToken.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", resp.InternalServerError("error signing refresh token").WithTracePrefix("error").AddTrace(err)
	}
	return tokenStr, nil
}

func parseAccessToken(tokenStr, secret string) (*models.AccessClaims, *resp.Response) {
	claims := &models.AccessClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, resp.Unauthorized("error parsing access token").WithTracePrefix("error").AddTrace(err)
	}
	if !token.Valid {
		return nil, resp.Unauthorized("invalid token")
	}
	return claims, nil
}

func parseRefreshToken(tokenStr, secret string) (*models.RefreshClaims, *resp.Response) {
	claims := &models.RefreshClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, resp.Unauthorized("error parsing refresh token").WithTracePrefix("error").AddTrace(err)
	}
	if !token.Valid {
		return nil, resp.Unauthorized("invalid token")
	}
	return claims, nil
}
