package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testRefresh(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("refresh@mail.com", ValidPassword).Register().Login()

	oldAccess := user.auth.AccessToken
	oldRefresh := user.auth.RefreshToken

	oldClient := client.WithAuth(&AuthContext{
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

		obj := resp.JSON().Object()
		obj.Value("error_id").String().IsEqual(string(apierr.TokenInvalid))
		obj.Value("message").String().IsEqual("refresh token is invalid")
	})

	t.Run("ConcurrentRefresh", func(t *testing.T) {
		// New user for this test
		concUser := client.WithCredentials("concurrent_refresh@mail.com", ValidPassword).Register().Login()
		refreshToken := concUser.auth.RefreshToken

		concurrency := 5
		results := make(chan int, concurrency)

		httpClient := &http.Client{Timeout: 5 * time.Second}

		var wg sync.WaitGroup
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
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

		go func() {
			wg.Wait()
			close(results)
		}()

		successCount := 0
		failCount := 0
		otherCount := 0

		for status := range results {
			if status == http.StatusOK {
				successCount++
			} else if status == http.StatusUnauthorized {
				failCount++
			} else {
				otherCount++
			}
		}

		require.Zero(t, otherCount, "Unexpected status code(s) returned: %d codes", otherCount)
		require.Equal(t, 1, successCount, "Only one refresh request should succeed")
		require.Equal(t, concurrency-1, failCount, "All other refresh requests should fail")
	})

	t.Run("TokenLinkage", func(t *testing.T) {
		// 1. Setup: User with two sessions
		client1 := suite.NewClient(t).WithCredentials("linkage1@mail.com", ValidPassword)
		user1 := client1.Register().Login()
		auth1 := user1.auth

		user2 := client1.Login()
		auth2 := user2.auth

		require.NotEqual(t, auth1.AccessToken, auth2.AccessToken)

		// 2. Attack: Use Access Token from Session 1 with Refresh Token from Session 2
		suite.NewClient(t).GET("/sessions/me").
			WithCookie("access_token", auth1.AccessToken).
			WithCookie("refresh_token", auth2.RefreshToken).
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.TokenMismatchDuringAuth).
			HasMessage("access token does not belong to this refresh token")

		// 3. Real linkage test: Old Access Token + New Refresh Token (Same Session)
		client3 := suite.NewClient(t).WithCredentials("linkage3@mail.com", ValidPassword)
		user3 := client3.Register().Login()
		oldAuth := *user3.auth

		// Refresh once
		refreshed := user3.Refresh()
		newAuth := refreshed.auth

		// Now we have:
		// S3: { Access-Old (valid for 15m), Access-New }
		// S3: { Refresh-Old (revoked), Refresh-New (valid) }

		// Using Access-Old with Refresh-New
		// Both have SessionID = S3.
		// Refresh-New was issued specifically to replace Access-Old.
		// If linkage is enforced, Refresh-New should only work with Access-New (or at least check AccessJTI).

		suite.NewClient(t).GET("/sessions/me").
			WithCookie("access_token", oldAuth.AccessToken).
			WithCookie("refresh_token", newAuth.RefreshToken).
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.TokenMismatchDuringAuth).
			HasMessage("access token does not belong to this refresh token")

		// 4. Positive case: New Access Token + New Refresh Token should work
		suite.NewClient(t).GET("/sessions/me").
			WithCookie("access_token", newAuth.AccessToken).
			WithCookie("refresh_token", newAuth.RefreshToken).
			Expect(http.StatusOK)
	})
}
