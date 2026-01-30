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
	t.Run("RefreshSuccess", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("refresh_success@mail.com", ValidPassword).Register().Login()

		oldAccess := user.auth.AccessToken
		oldRefresh := user.auth.RefreshToken

		if user.auth == nil {
			t.Errorf("Refresh called on unauthenticated client")
		}

		resp := user.POST("/auth/refresh").
			WithCookie("refresh_token", user.auth.RefreshToken).
			Expect(http.StatusOK)

		access := resp.resp.Cookie("access_token")
		refresh := resp.resp.Cookie("refresh_token")

		if access.Raw() == nil {
			t.Errorf("Expected access cookie after refresh but got nil")
		}

		if refresh.Raw() == nil {
			t.Errorf("Expected refresh cookie after refresh but got nil")
		}

		refreshed := user.WithAuth(&AuthContext{
			AccessToken:  access.Value().Raw(),
			RefreshToken: refresh.Value().Raw(),
		})

		require.NotEqual(t, oldAccess, refreshed.auth.AccessToken, "Access token should change after refresh")
		require.NotEqual(t, oldRefresh, refreshed.auth.RefreshToken, "Refresh token should change after refresh")
	})

	t.Run("UseOldTokenError", func(t *testing.T) {
		client := suite.NewClient(t)
		oldAuth := client.WithCredentials("old_token_user@mail.com", ValidPassword).Register().Login()

		oldAuth.Refresh() // Refresh the token, invalidating the old one, never save the new authed user

		// Now try to use the old access token
		client.WithAuth(oldAuth.auth).GET("/sessions").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.TokenReuseIdentified).
			HasMessage("refresh token reuse not allowed")
	})

	t.Run("RefreshRevokedToken", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("revoked_refresh_user@mail.com", ValidPassword).Register().Login()
		refreshToken := user.auth.RefreshToken

		// Manually revoke the session in the database
		_, err := suite.DB.Exec(`
			UPDATE sessions 
			SET revoked_at = NOW() 
			WHERE token_id = (
				SELECT token_id 
				FROM sessions 
				WHERE identity_id = (
					SELECT id 
					FROM identities 
					WHERE entity_id = (
						SELECT id 
						FROM users 
						WHERE email = 'revoked_refresh_user@mail.com'
					)
				) LIMIT 1
			)
		`)
		require.NoError(t, err)

		// Attempt to use the revoked refresh token
		deniedClient := suite.NewClient(t)
		deniedClient.POST("/auth/refresh").
			WithCookie("refresh_token", refreshToken).
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SessionNotFound).
			HasMessage("session not found or revoked")
	})

	t.Run("ConcurrentRefresh", func(t *testing.T) {
		concUser := suite.NewClient(t).WithCredentials("concurrent_refresh@mail.com", ValidPassword).Register().Login()
		refreshToken := concUser.auth.RefreshToken

		concurrency := 5
		results := make(chan int, concurrency)
		httpClient := &http.Client{Timeout: 10 * time.Second}
		var wg sync.WaitGroup
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				req, _ := http.NewRequest("POST", suite.Server.URL+"/auth/refresh", nil)
				req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
				resp, err := httpClient.Do(req)
				if err != nil {
					results <- 0
					return
				}
				defer resp.Body.Close()
				results <- resp.StatusCode
			}()
		}

		wg.Wait()
		close(results)

		successCount := 0
		failCount := 0
		for status := range results {
			if status == http.StatusOK {
				successCount++
			} else if status == http.StatusUnauthorized {
				failCount++
			}
		}

		require.Equal(t, 1, successCount, "Only one refresh request should succeed")
		require.Equal(t, concurrency-1, failCount, "All other refresh requests should fail due to token reuse")
	})
}
