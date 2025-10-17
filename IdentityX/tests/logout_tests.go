package testing

import (
	"net/http"
	"testing"
)

func logoutNoTokens() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/logout").
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("missing refresh_token cookie")

		obj.Value("code").Number().Equal(401)
	}
}

func logoutNoRefresh(user *rllCtx) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/logout").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("missing refresh_token cookie")

		obj.Value("code").Number().Equal(401)
	}
}

func logoutSuccess(user *rllCtx) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/logout").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("Logged out")

		obj.Value("code").Number().Equal(200)
	}
}

func loggedOutAlready(user *rllCtx) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/logout").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("user already logged out")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("token already blacklisted")

		obj.Value("code").Number().Equal(400)
	}
}
