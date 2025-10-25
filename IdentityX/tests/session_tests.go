package testing

import (
	"net/http"
	"testing"
)

func listXSessions(user *accountContext, session_amount int) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.GET("/sessions").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")

		data := obj.Value("data").Array()
		data.Length().Equal(session_amount)
		user.sessionID = data.Element(session_amount - 1).Object().Value("session_id").String().Raw()
		user.sessionJIT = data.Element(session_amount - 1).Object().Value("token_id").String().Raw()

		obj.Value("code").Number().Equal(200)
	}
}

func listSessionsFail(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.GET("/sessions").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().Equal("AuthMW")
		obj.Value("message").String().Equal("refresh token is invalidated")

		obj.Value("code").Number().Equal(401)
	}
}

func revokeSessionByIDFail(user *accountContext) func(t *testing.T) {
  return func(t *testing.T) {
    e := createExpect(t)

		obj := e.DELETE("/sessions/" + user.sessionID).
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("can't revoke a currently active session, please logout instead")
		obj.Value("code").Number().Equal(400)
  }
}

func revokeSessionByIDSuccess(user *accountContext) func(t *testing.T) {
  return func(t *testing.T) {
    e := createExpect(t)

		obj := e.DELETE("/sessions/" + user.sessionID).
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("revoked session")
		obj.Value("code").Number().Equal(200)
  }
}

func revokeOtherSessions(user *accountContext) func(t *testing.T) {
  return func(t *testing.T) {
    e := createExpect(t)

		obj := e.DELETE("/sessions/others").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("revoked sessions")
		obj.Value("code").Number().Equal(200)
  }
}

func revokeAllSessions(user *accountContext) func(t *testing.T) {
  return func(t *testing.T) {
    e := createExpect(t)

		obj := e.DELETE("/sessions").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("revoked sessions")
		obj.Value("code").Number().Equal(200)
  }
}
