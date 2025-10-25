package testing

import (
	"net/http"
	"testing"
)

func ping() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/ping/public").
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("pong")

		obj.Value("code").Number().Equal(200)
	}
}

func privatePingFailure() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/ping/private").
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().Equal("AuthMW")
		obj.Value("message").String().Equal("missing access_token cookie")

		obj.Value("code").Number().Equal(401)
	}
}

func privatePingSuccess(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/ping/private").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()


		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Contains("pong to you")

		obj.Value("code").Number().Equal(200)
	}
}
