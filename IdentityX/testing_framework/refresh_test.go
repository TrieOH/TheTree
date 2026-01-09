package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testRefresh(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.NewUser("refresh@mail.com", ValidPassword).Register().Login()

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
			HasErrID(apierr.SessionUnauthorized).
			HasMessage("session not found or revoked")
	})

	t.Run("RefreshRevokedToken", func(t *testing.T) {
		deniedClient := suite.NewClient(t)

		resp := deniedClient.expect.POST("/auth/refresh").
			WithCookie("refresh_token", oldRefresh).
			Expect().
			Status(http.StatusUnauthorized)

		resp.JSON().Object().Value("error_id").String().IsEqual(string(apierr.TokenInvalid))
		resp.JSON().Object().Value("message").String().IsEqual("refresh token is invalid")
	})

	t.Run("ConcurrentRefresh", func(t *testing.T) {
		// New user for this test
		concUser := client.NewUser("concurrent_refresh@mail.com", ValidPassword).Register().Login()
		refreshToken := concUser.auth.RefreshToken

		concurrency := 5
		results := make(chan int, concurrency)

		httpClient := &http.Client{Timeout: 5 * time.Second}

		for i := 0; i < concurrency; i++ {
			go func() {
				req, err := http.NewRequest("POST", suite.Server.URL+"/auth/refresh", nil)
				if err != nil {
					results <- 0
					return
				}
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
				resp, err := httpClient.Do(req)
				if err != nil {
					// Treat network error as failure to refresh (or 0)
					results <- 0
					return
				}
				defer resp.Body.Close()
				results <- resp.StatusCode
			}()
		}

		successCount := 0
		failCount := 0
		otherCount := 0

		for i := 0; i < concurrency; i++ {
			status := <-results
			if status == http.StatusOK {
				successCount++
			} else if status == http.StatusUnauthorized {
				failCount++
			} else {
				otherCount++
			}
		}

		require.Zero(t, otherCount, "Unexpected status code(s) returned: %v codes", otherCount)
		require.Equal(t, 1, successCount, "Only one refresh request should succeed")
		require.Equal(t, concurrency-1, failCount, "All other refresh requests should fail")
	})
}
