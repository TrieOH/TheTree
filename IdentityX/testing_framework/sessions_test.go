package testing

import (
	"net/http"
	"testing"
)

func testSessions(t *testing.T, suite *TestSuite) {
	// Create user in parent test context
	client := suite.Client(t)
	user := client.User("sessions@mail.com", ValidPassword).Register().Login()

	t.Run("ListSessions", func(t *testing.T) {
		// Create a new client with subtest's t for the authenticated client
		// We need to preserve the auth context but use the new t
		authClient := suite.Client(t).Auth(user.auth)
		arr := authClient.GET("/sessions").
			Expect(http.StatusOK).
			DataArray()

		arr.Length().IsEqual(1)
	})

	t.Run("MultipleLoginsSessions", func(t *testing.T) {
		// Create 3 more sessions (we already have 1)
		client2 := suite.Client(t)
		client2.User(user.Email, user.Password).Login()

		client3 := suite.Client(t)
		client3.User(user.Email, user.Password).Login()

		client4 := suite.Client(t)
		user4 := client4.User(user.Email, user.Password).Login()

		arr := user4.AuthedClient().GET("/sessions").
			Expect(http.StatusOK).
			DataArray()

		arr.Length().IsEqual(4)

		// Get oldest session ID to revoke
		currentSessionID := arr.Value(0).Object().Value("session_id").String().Raw()
		oldestSessionID := arr.Value(3).Object().Value("session_id").String().Raw()

		// Can't revoke current session
		user4.AuthedClient().DELETE("/sessions/"+currentSessionID).
			Expect(http.StatusForbidden).
			Error("go-auth-test", "cannot revoke the currently active session")

		// Revoke first session
		user4.AuthedClient().DELETE("/sessions/"+oldestSessionID).
			Expect(http.StatusOK).
			Success("go-auth-test", "revoked session")

		// Should have 3 sessions now
		user4.AuthedClient().GET("/sessions").
			Expect(http.StatusOK).
			DataArray().
			Length().IsEqual(3)
	})

	t.Run("RevokeOtherSessions", func(t *testing.T) {
		client := suite.Client(t)
		user := client.User("revoke-others@mail.com", ValidPassword).Register().Login()

		// Create more sessions
		suite.Client(t).User(user.Email, user.Password).Login()
		suite.Client(t).User(user.Email, user.Password).Login()

		user.AuthedClient().DELETE("/sessions/others").
			Expect(http.StatusOK).
			Success("go-auth-test", "revoked sessions")

		// Should only have current session
		user.AuthedClient().GET("/sessions").
			Expect(http.StatusOK).
			DataArray().
			Length().IsEqual(1)
	})

	t.Run("SessionInfo", func(t *testing.T) {
		client := suite.Client(t)
		user := client.User("session-me@mail.com", ValidPassword).Register().Login()

		data := user.AuthedClient().GET("/sessions/me").
			Expect(http.StatusOK).
			Value()

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
		client := suite.Client(t)
		user := client.User("revoke-all@mail.com", ValidPassword).Register().Login()

		// Create more sessions
		suite.Client(t).User(user.Email, user.Password).Login()
		suite.Client(t).User(user.Email, user.Password).Login()

		user.AuthedClient().DELETE("/sessions").
			Expect(http.StatusOK).
			Success("go-auth-test", "revoked sessions")

		// Session should be invalid
		user.AuthedClient().GET("/sessions").
			Expect(http.StatusUnauthorized).
			Error("AuthMW", "refresh token is revoked")
	})
}
