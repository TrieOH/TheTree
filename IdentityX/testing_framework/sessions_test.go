package testing

import (
	"GoAuth/internal/errx"
	"context"
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
			HasErrID(errx.SessionSelfRevokeForbidden).
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
			HasErrID(errx.SessionRevoked).
			HasMessage("session not found or revoked")
	})

	// ExpiredSessionNotListed tests that expired sessions are not returned when listing active sessions.
	// This ensures that only currently valid sessions are visible to the user.
	// The test manually inserts an expired session directly into the database for a registered user.
	// It then attempts to list sessions for that user, expecting only the active session (created by login)
	// to be returned, and not the manually inserted expired one.
	t.Run("ExpiredSessionNotListed", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("expired@mail.com", ValidPassword).Register().Login()

		ctx := context.Background()

		// Manually insert an expired session for this user
		var userID string
		err := suite.DB.QueryRow(ctx, "SELECT id FROM users WHERE email = 'expired@mail.com'").Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to get user ID: %v", err)
		}

		var identityID string
		err = suite.DB.QueryRow(ctx, `
			INSERT INTO identities (type, entity_id)
			VALUES ('client', $1)
			ON CONFLICT (type, entity_id) DO UPDATE SET type = 'client'
			RETURNING id
		`, userID).Scan(&identityID)
		if err != nil {
			t.Fatalf("Failed to create session identity: %v", err)
		}

		_, err = suite.DB.Exec(ctx, `
			INSERT INTO sessions (
				identity_id, issued_at, user_agent, user_ip, expires_at, created_at, updated_at, user_type
			) VALUES (
				$1,
				NOW() - INTERVAL '2 days',
				'Expired Agent',
				'127.0.0.1',
				NOW() - INTERVAL '1 day',
				NOW(),
				NOW(),
				'client'
			)
		`, identityID)
		if err != nil {
			t.Fatalf("Failed to insert expired session: %v", err)
		}

		// Verify that the expired session is NOT in the list
		// Should only have the active login session (1), ignoring the manually inserted expired one
		user.GET("/sessions").
			Expect(http.StatusOK).
			RequireDataArray().Length().IsEqual(1)
	})

	// RevokedSessionNotListed tests that manually revoked sessions are not returned when listing active sessions.
	// This verifies that sessions explicitly marked as revoked are correctly filtered out.
	// The test manually inserts a revoked session directly into the database for a registered user.
	// It then lists sessions for that user, expecting only the active session (created by login)
	// to be returned, and not the manually inserted revoked one.
	t.Run("RevokedSessionNotListed", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("revoked@mail.com", ValidPassword).Register().Login()

		ctx := context.Background()

		// Manually insert a revoked session for this user
		var userID string
		err := suite.DB.QueryRow(ctx, "SELECT id FROM users WHERE email = 'revoked@mail.com'").Scan(&userID)
		if err != nil {
			t.Fatalf("Failed to get user ID: %v", err)
		}

		var identityID string
		err = suite.DB.QueryRow(ctx, `
			INSERT INTO identities (type, entity_id)
			VALUES ('client', $1)
			ON CONFLICT (type, entity_id) DO UPDATE SET type = 'client'
			RETURNING id
		`, userID).Scan(&identityID)
		if err != nil {
			t.Fatalf("Failed to create session identity: %v", err)
		}

		_, err = suite.DB.Exec(ctx, `
			INSERT INTO sessions (
				identity_id, issued_at, user_agent, user_ip, revoked_at, created_at, updated_at, expires_at, user_type
			) VALUES (
				$1,
				NOW() - INTERVAL '2 days',
				'Expired Agent',
				'127.0.0.1',
				NOW() - INTERVAL '1 day',
				NOW(),
				NOW(),
				NOW() + INTERVAL '1 day',
				'client'
			)
		`, identityID)
		if err != nil {
			t.Fatalf("Failed to insert revoked session: %v", err)
		}

		// Verify that the revoked session is NOT in the list
		// Should only have the active login session (1), ignoring the manually inserted revoked one
		user.GET("/sessions").
			Expect(http.StatusOK).
			RequireDataArray().Length().IsEqual(1)
	})

	// SessionLeakage tests that session isolation is correctly enforced between different identity types (client vs. project user).
	// This is a critical security test to prevent one user type from accessing or even seeing sessions belonging to another type,
	// even if both identity types might logically be linked to the same underlying individual or email.
	// The test creates a client user and a project user (both potentially using the same email).
	// It then logs in as each user, creating separate sessions for each identity type.
	// Finally, it attempts to list sessions while authenticated as the client user,
	// expecting to see only the client user's session and ensuring that the project user's session is not exposed.
	t.Run("SessionLeakage", func(t *testing.T) {
		// This test ensures that a user with a 'client' identity cannot see sessions
		// from a 'project' identity, even if they are logically the same user.

		// 1. Create a client user
		clientUser := suite.NewClient(t).WithCredentials("leakage-client@mail.com", ValidPassword)
		clientUser.Register()

		// 2. Create a project and a project user with the same email
		projectOwner := suite.NewClient(t).WithCredentials("leakage-project-owner@mail.com", ValidPassword)
		project := projectOwner.Register().Login().CreateProject("Leakage Test Project")

		projectUser := suite.NewClient(t).WithCredentials("leakage-client@mail.com", ValidPassword)
		projectUser.ProjectRegister(project.projectID)

		// 3. Log in as both the client and project user to create two separate sessions
		clientSession := clientUser.Login()
		projectUser.ProjectLogin(project.projectID) // This creates a session for the project user

		// 4. As the client user, list the sessions
		// We expect to only see the client user's session, not the project user's.
		sessions := clientSession.GET("/sessions").
			Expect(http.StatusOK).
			RequireDataArray()

		// 5. Verify that only one session is returned
		sessions.Length().IsEqual(1)

		// Optional: Verify the user_type of the returned session
		sessionType := sessions.Value(0).Object().Value("user_type").String().Raw()
		if sessionType != "client" {
			t.Errorf("Expected user_type to be 'client', but got '%s'", sessionType)
		}
	})
}
