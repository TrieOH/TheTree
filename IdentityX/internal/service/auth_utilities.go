package service

import (
	"time"

	"GoAuth/internal/models"
	"GoAuth/internal/repository"

	resp "github.com/MintzyG/GoResponse/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
)

func newAccessToken(dbUser repository.User) (string, *resp.Response) {
	var userSub models.UserSubJWT
	if err := copier.Copy(&userSub, &dbUser); err != nil {
		return "", resp.InternalServerError("errir creating user sub for access token").WithTracePrefix("error").AddTrace(err)
	}

	expiration := time.Now().Add(1 * time.Hour)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "go_auth_server",
			"sub": userSub,
			"exp": expiration.Unix(),
		})
	tokenStr, err := accessToken.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", resp.InternalServerError("error signing access token").WithTracePrefix("error").AddTrace(err)
	}
	return tokenStr, nil
}

func newRefreshToken() (string, *resp.Response) {
	expiration := time.Now().Add(7 * 24 * time.Hour)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "go_auth_server",
			"sub": "refreshMetaData",
			"exp": expiration.Unix(),
		})
	tokenStr, err := refreshToken.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", resp.InternalServerError("error signing refresh token").WithTracePrefix("error").AddTrace(err)
	}
	return tokenStr, nil
}
