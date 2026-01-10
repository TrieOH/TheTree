package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"
)

func testSessions(t *testing.T, suite *TestSuite) {
	// Create user in parent test context
	client := suite.NewClient(t)
	user := client.WithCredentials("sessions@mail.com", ValidPassword).Register().Login()

	t.Run("ListSessions", func(t *testing.T) {
		// Create a new client with subtest's t for the authenticated client
		// We need to preserve the auth context but use the new t
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.GET("/sessions").
			Expect(http.StatusOK).
			RequireDataArray().Length().IsEqual(1)
	})

	t.Run("MultipleLoginsSessions", func(t *testing.T) {
		// Create 3 more sessions (we already have 1)
		client2 := suite.NewClient(t)
		client2.WithCredentials(user.email, user.password).Login()

		client3 := suite.NewClient(t)
		client3.WithCredentials(user.email, user.password).Login()

		client4 := suite.NewClient(t)
		user4 := client4.WithCredentials(user.email, user.password).Login()

		arr := user4.GET("/sessions").
			Expect(http.StatusOK).
			RequireDataArray()

		arr.Length().IsEqual(4)

		// Get oldest session ID to revoke
		currentSessionID := arr.Value(0).Object().Value("session_id").String().Raw()
		oldestSessionID := arr.Value(3).Object().Value("session_id").String().Raw()

		// Can't revoke current session
		user4.DELETE("/sessions/" + currentSessionID).
			Expect(http.StatusForbidden).
			HasErrID(apierr.SessionSelfRevokeForbidden).
			HasMessage("cannot revoke the currently active session")

		// Revoke first session
		user4.DELETE("/sessions/" + oldestSessionID).
			Expect(http.StatusOK).
			HasMessage("revoked session")

		// Should have 3 sessions now
		user4.GET("/sessions").
			Expect(http.StatusOK).
			RequireDataArray().Length().IsEqual(3)
	})

	t.Run("RevokeOtherSessions", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("revoke-others@mail.com", ValidPassword).Register().Login()

		// Create more sessions
		suite.NewClient(t).WithCredentials(user.email, user.password).Login()
		suite.NewClient(t).WithCredentials(user.email, user.password).Login()

		user.DELETE("/sessions/others").
			Expect(http.StatusOK).
			HasMessage("revoked sessions")

		// Should only have current session
		user.GET("/sessions").
			Expect(http.StatusOK).
			RequireDataArray().Length().IsEqual(1)
	})

	t.Run("SessionInfo", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("session-me@mail.com", ValidPassword).Register().Login()

		data := user.GET("/sessions/me").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"refresh_expire_date": AnyNumber{},
			"access": map[string]interface{}{
				"iss": "GoAuth",
				"exp": AnyNumber{},
				"iat": AnyNumber{},
				"jti": AnyUUID{},
				"sub": map[string]interface{}{
					"id":         AnyUUID{},
					"email":      "session-me@mail.com",
					"project_id": nil,
					"user_type":  "client",
					"metadata":   nil,
					"session_id": AnyUUID{},
					"user_agent": AnyString{},
					"user_ip":    AnyString{},
				},
			},
		}

		Validate(t, data, spec)
	})

	t.Run("RevokeAllSessions", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("revoke-all@mail.com", ValidPassword).Register().Login()

		// Create more sessions
		suite.NewClient(t).WithCredentials(user.email, user.password).Login()
		suite.NewClient(t).WithCredentials(user.email, user.password).Login()

		user.DELETE("/sessions").
			Expect(http.StatusOK).
			HasMessage("revoked sessions")

		// Session should be invalid
		user.GET("/sessions").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SessionUnauthorized).
			HasMessage("session not found or revoked")
	})

	t.Run("ExpiredSessionNotListed", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("expired@mail.com", ValidPassword).Register().Login()

		// Manually insert an expired session for this user
		_, err := suite.DB.Exec(`
			INSERT INTO sessions (
				user_id, issued_at, user_agent, user_ip, expires_at, user_type, created_at, updated_at
			) VALUES (
				(SELECT id FROM users WHERE email = 'expired@mail.com'),
				NOW() - INTERVAL '2 days',
				'Expired Agent',
				'127.0.0.1',
				NOW() - INTERVAL '1 day',
				'client',
				NOW(),
				NOW()
			)
		`)
		if err != nil {
			t.Fatalf("Failed to insert expired session: %v", err)
		}

		// Verify that the expired session is NOT in the list
		// Should only have the active login session (1), ignoring the manually inserted expired one
		user.GET("/sessions").
			Expect(http.StatusOK).
			RequireDataArray().Length().IsEqual(1)
	})

	t.Run("RevokedSessionNotListed", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("revoked@mail.com", ValidPassword).Register().Login()

		// Manually insert an expired session for this user
		_, err := suite.DB.Exec(`
			INSERT INTO sessions (
				user_id, issued_at, user_agent, user_ip, revoked_at, user_type, created_at, updated_at, expires_at
			) VALUES (
				(SELECT id FROM users WHERE email = 'expired@mail.com'),
				NOW() - INTERVAL '2 days',
				'Expired Agent',
				'127.0.0.1',
				NOW() - INTERVAL '1 day',
				'client',
				NOW(),
				NOW(),
				NOW() + INTERVAL '1 day'
			)
		`)
		if err != nil {
			t.Fatalf("Failed to insert expired session: %v", err)
		}

		// Verify that the revoked session is NOT in the list
		// Should only have the active login session (1), ignoring the manually inserted revoked one
		user.GET("/sessions").
			Expect(http.StatusOK).
			RequireDataArray().Length().IsEqual(1)
	})
}
