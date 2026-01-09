package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func testRefresh(t *testing.T, suite *TestSuite) {
	client := suite.Client(t)
	user := client.User("refresh@mail.com", ValidPassword).Register().Login()

	oldAccess := user.auth.AccessToken
	oldRefresh := user.auth.RefreshToken

	oldClient := client.Auth(&AuthContext{
		AccessToken:  oldAccess,
		RefreshToken: oldRefresh,
	})

	t.Run("RefreshSuccess", func(t *testing.T) {
		refreshed := user.WithT(t).Refresh()

		require.NotEqual(t, oldAccess, refreshed.auth.AccessToken, "Access token should change after refresh")
		require.NotEqual(t, oldRefresh, refreshed.auth.RefreshToken, "Refresh token should change after refresh")
	})

	t.Run("UseOldTokenError", func(t *testing.T) {
		oldClient.WithT(t).GET("/sessions").
			Expect(http.StatusUnauthorized).
			ExpectErrorID(apierr.TokenRevoked).
			MessageContains("refresh token is revoked")
	})

	t.Run("RefreshRevokedToken", func(t *testing.T) {
		oldClient := suite.Client(t)

		resp := oldClient.expect.POST("/auth/refresh").
			WithCookie("refresh_token", oldRefresh).
			Expect().
			Status(http.StatusUnauthorized)

		resp.JSON().Object().Value("error_id").String().IsEqual(string(apierr.TokenRevoked))
		resp.JSON().Object().Value("message").String().IsEqual("refresh token revoked")
	})
}
