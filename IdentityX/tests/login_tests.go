package testing

import (
	"net/http"
	"testing"
)

func loginWrongPassword(user *rllCtx) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/login").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": user.SuccessEmail,
				"password": "123",
			}).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("invalid email or password")

		obj.Value("code").Number().Equal(401)
	}
}

func loginWrongEmail(user *rllCtx) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/login").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "wrong@email.com",
				"password": user.SuccessPasword, 
			}).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("invalid email or password")

		obj.Value("code").Number().Equal(401)
	}
}

func LoginWrongEmailAndPasword() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/login").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "wrong@email.com",
				"password": "Wr0ngP4$$",
			}).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("invalid email or password")

		obj.Value("code").Number().Equal(401)
	}
}

func loginSuccess(user *rllCtx) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		resp := e.POST("/auth/login").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    user.SuccessEmail,
				"password": user.SuccessPasword,
			}).
			Expect().
			Status(http.StatusOK)

		obj := resp.JSON().Object()
		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("Logged in")
		obj.Value("code").Number().Equal(200)

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

		t.Logf("Access token: %s", user.accessToken)
		t.Logf("Refresh token: %s", user.refreshToken)
	}
}
