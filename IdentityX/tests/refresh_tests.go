package testing

import (
	"net/http"
	"testing"
)

func refreshTokensSuccess(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		resp := e.POST("/auth/refresh").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK)

		obj := resp.JSON().Object()
		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("Refreshed tokens")

		access := resp.Cookie("access_token")
		if access == nil || access.Raw() == nil {
			t.Fatalf("expected access_token cookie, got nil")
		}

		val := access.Value().Raw()
		if val == "" {
			t.Fatalf("access_token cookie value is empty")
		}
		user.accessToken = val

		refresh := resp.Cookie("refresh_token")
		if refresh == nil || refresh.Raw() == nil {
			t.Fatalf("expected refresh_token cookie, got nil")
		}

		val = refresh.Value().Raw()
		if val == "" {
			t.Fatalf("refresh_token cookie value is empty")
		}
		user.refreshToken = val

		obj.Value("code").Number().Equal(200)
	}
}
