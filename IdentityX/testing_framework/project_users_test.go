package testing

import (
	"GoAuth/internal/apierr"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func testProjectUsers(t *testing.T, suite *TestSuite) {
	// Create user in parent test context
	client := suite.NewClient(t)
	user := client.NewUser("client@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("test project")

	// Create with the same email as client to prove it's in a different environment
	projectUser := client.NewUser("client@mail.com", ValidPassword).
		ProjectRegister(user.ProjectID).
		ProjectLogin(user.ProjectID)

	t.Run("ProjectUsersRegister", func(t *testing.T) {
		var wrongFormatID = "wrong-format"
		t.Run("WrongFormatIDProjectRegister", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + wrongFormatID + "/register").
				WithBody(map[string]interface{}{
					"email":    user.Email,
					"password": user.Password,
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.ProjectInvalidID).
				HasMessage("invalid project id")
		})

		var nonexistentID = "d917a262-199a-41b1-af43-930ea8da1c75"
		t.Run("InvalidProjectRegister", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + nonexistentID + "/register").
				WithBody(map[string]interface{}{
					"email":    user.Email,
					"password": user.Password,
				}).
				Expect(http.StatusBadRequest).
				HasMessage("invalid reference").
				TraceContains("violates foreign key costraint")
		})

		t.Run("ValidationProjectRegister", func(t *testing.T) {
			for _, spec := range ValidationTests {
				spec := spec
				t.Run(spec.Name, func(t *testing.T) {
					client := suite.NewClient(t)
					client.POST("/projects/" + user.ProjectID + "/register").
						WithBody(map[string]interface{}{
							"email":    spec.Email,
							"password": spec.Pass,
						}).
						Expect(http.StatusBadRequest).
						HasErrID(apierr.RequestValidationError).
						ValidationError(spec.Errors...)
				})
			}
		})

		t.Run("ValidationProjectRegister", func(t *testing.T) {
			for i, spec := range WeakPasswordTests {
				spec := spec
				i := i
				t.Run(spec.Name, func(t *testing.T) {
					client := suite.NewClient(t)
					client.POST("/projects/" + user.ProjectID + "/register").
						WithBody(map[string]interface{}{
							"email":    fmt.Sprintf("weak%d@mail.com", i),
							"password": spec.Password,
						}).
						Expect(http.StatusBadRequest).
						HasErrID(apierr.RequestValidationError).
						ValidationError(spec.Errors...)
				})
			}
		})

		t.Run("DuplicateEmailProjectRegister", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + user.ProjectID + "/register").
				WithBody(map[string]interface{}{
					"email":    user.Email,
					"password": user.Password,
				}).
				Expect(http.StatusConflict).
				HasErrID(apierr.AuthEmailAlreadyUsed).
				HasMessage("error registering user").
				TraceContains("email already in use")
		})

		t.Run("InvalidSchemaTypeRegister", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/"+user.ProjectID+"/register").
				WithQuery("schema_type", "invalid").
				WithBody(map[string]interface{}{
					"email":    "invalid_schema@email.com",
					"password": user.Password,
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.SchemaInvalidSchemaType).
				HasMessage("invalid schema type")
		})

		t.Run("FlowIDSameAsTypeRegister", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/"+user.ProjectID+"/register").
				WithQuery("schema_type", "context").
				WithQuery("flow_id", "context").
				WithBody(map[string]interface{}{
					"email":    "flow_same_as_type@email.com",
					"password": user.Password,
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.SchemaInvalidFlowID).
				HasMessage("flow id can't be the same as a schema type")
		})

		t.Run("MetadataRegisterOnCoreDenied", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + user.ProjectID + "/register").
				WithBody(map[string]interface{}{
					"email":    "metadata_denied@email.com",
					"password": user.Password,
					"custom_fields": map[string]interface{}{
						"curso": "Ciência da Computação",
					},
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.SchemaMetadataNotAllowed).
				HasMessage("custom fields are not allowed for core schema")
		})

		t.Run("SuccessProjectRegister", func(t *testing.T) {
			client := suite.NewClient(t)
			client.NewUser("new@mail.com", ValidPassword).ProjectRegister(user.ProjectID)
		})
	})

	t.Run("ProjectUsersLogin", func(t *testing.T) {
		t.Run("WrongPassword", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + user.ProjectID + "/login").
				WithBody(map[string]string{
					"email":    projectUser.Email,
					"password": "WrongPass123!",
				}).
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.AuthInvalidCredentials).
				HasMessage("invalid email or password")
		})

		t.Run("WrongEmail", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + user.ProjectID + "/login").
				WithBody(map[string]string{
					"email":    "wrong@mail.com",
					"password": projectUser.Password,
				}).
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.AuthInvalidCredentials).
				HasMessage("invalid email or password")
		})

		t.Run("Success", func(t *testing.T) {
			client := suite.NewClient(t)
			client.NewUser(projectUser.Email, projectUser.Password).ProjectLogin(user.ProjectID)
		})

		t.Run("Logout", func(t *testing.T) {
			client := suite.NewClient(t)
			loggedInUser := client.NewUser(projectUser.Email, projectUser.Password).ProjectLogin(user.ProjectID)
			loggedInUser.Logout()

			// Try using revoked session
			loggedInUser.AuthedClient().POST("/auth/logout").
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.SessionUnauthorized).
				HasMessage("session not found or revoked")
		})
	})

	sessionUser := client.NewUser("sessions@mail.com", ValidPassword).
		ProjectRegister(user.ProjectID).
		ProjectLogin(user.ProjectID)

	t.Run("ProjectUsersSession", func(t *testing.T) {
		t.Run("ListSessions", func(t *testing.T) {
			// Create a new client with subtest's t for the authenticated client
			// We need to preserve the auth context but use the new t
			authClient := suite.NewClient(t).Auth(sessionUser.auth)
			authClient.GET("/sessions").
				Expect(http.StatusOK).
				RequireDataArray().Length().IsEqual(1)
		})

		t.Run("MultipleLoginsSessions", func(t *testing.T) {
			// Create 3 more sessions (we already have 1)
			client2 := suite.NewClient(t)
			client2.NewUser(sessionUser.Email, sessionUser.Password).ProjectLogin(user.ProjectID)

			client3 := suite.NewClient(t)
			client3.NewUser(sessionUser.Email, sessionUser.Password).ProjectLogin(user.ProjectID)

			client4 := suite.NewClient(t)
			user4 := client4.NewUser(sessionUser.Email, sessionUser.Password).ProjectLogin(user.ProjectID)

			arr := user4.AuthedClient().GET("/sessions").
				Expect(http.StatusOK).
				RequireDataArray()

			arr.Length().IsEqual(4)

			// Get oldest session ID to revoke
			currentSessionID := arr.Value(0).Object().Value("session_id").String().Raw()
			oldestSessionID := arr.Value(3).Object().Value("session_id").String().Raw()

			// Can't revoke current session
			user4.AuthedClient().DELETE("/sessions/" + currentSessionID).
				Expect(http.StatusForbidden).
				HasErrID(apierr.SessionSelfRevokeForbidden).
				HasMessage("cannot revoke the currently active session")

			// Revoke first session
			user4.AuthedClient().DELETE("/sessions/" + oldestSessionID).
				Expect(http.StatusOK).
				HasErrID(apierr.SessionRevoked).
				HasMessage("revoked session")

			// Should have 3 sessions now
			user4.AuthedClient().GET("/sessions").
				Expect(http.StatusOK).
				RequireDataArray().Length().IsEqual(3)
		})

		t.Run("RevokeOtherSessions", func(t *testing.T) {
			client := suite.NewClient(t)
			revokeOthers := client.NewUser("revoke-others-project@mail.com", ValidPassword).
				ProjectRegister(user.ProjectID).
				ProjectLogin(user.ProjectID)

			// Create more sessions
			suite.NewClient(t).NewUser(revokeOthers.Email, revokeOthers.Password).ProjectLogin(user.ProjectID)
			suite.NewClient(t).NewUser(revokeOthers.Email, revokeOthers.Password).ProjectLogin(user.ProjectID)

			revokeOthers.AuthedClient().DELETE("/sessions/others").
				Expect(http.StatusOK).
				HasMessage("revoked sessions")

			// Should only have current session
			revokeOthers.AuthedClient().GET("/sessions").
				Expect(http.StatusOK).
				RequireDataArray().Length().IsEqual(1)
		})

		t.Run("SessionInfo", func(t *testing.T) {
			client := suite.NewClient(t)
			infoUser := client.NewUser("session-me@mail.com", ValidPassword).
				ProjectRegister(user.ProjectID).
				ProjectLogin(user.ProjectID)

			data := infoUser.AuthedClient().GET("/sessions/me").
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
						"project_id": AsString{user.ProjectID, AnyUUID{}},
						"user_type":  "project",
						"metadata":   map[string]interface{}{},
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
			revoked := client.NewUser("revoke-all@mail.com", ValidPassword).
				ProjectRegister(user.ProjectID).
				ProjectLogin(user.ProjectID)

			// Create more sessions
			suite.NewClient(t).NewUser(revoked.Email, revoked.Password).ProjectLogin(user.ProjectID)
			suite.NewClient(t).NewUser(revoked.Email, revoked.Password).ProjectLogin(user.ProjectID)

			revoked.AuthedClient().DELETE("/sessions").
				Expect(http.StatusOK).
				HasMessage("revoked sessions")

			// Session should be invalid
			revoked.AuthedClient().GET("/sessions").
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.SessionUnauthorized).
				HasMessage("session not found or revoked")
		})
	})

	t.Run("ProjectUserRefresh", func(t *testing.T) {
		client := suite.NewClient(t)
		refreshUSer := client.NewUser("refresh@mail.com", ValidPassword).
			ProjectRegister(user.ProjectID).
			ProjectLogin(user.ProjectID)

		oldAccess := refreshUSer.auth.AccessToken
		oldRefresh := refreshUSer.auth.RefreshToken

		refreshUSer.Refresh()

		require.NotEqual(t, oldAccess, refreshUSer.auth.AccessToken, "Access token should change after refresh")
		require.NotEqual(t, oldRefresh, refreshUSer.auth.RefreshToken, "Refresh token should change after refresh")

		// Old tokens should be invalid
		oldClient := client.Auth(&AuthContext{
			AccessToken:  oldAccess,
			RefreshToken: oldRefresh,
		})

		oldClient.GET("/sessions").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SessionUnauthorized).
			HasMessage("session not found or revoked")
	})

	t.Run("ProjectUsersProjects", func(t *testing.T) {
		client := suite.NewClient(t)
		nested := client.NewUser("nested_creator@mail.com", ValidPassword).
			ProjectRegister(user.ProjectID).
			ProjectLogin(user.ProjectID)

		t.Run("CreateProject", func(t *testing.T) {
			authClient := suite.NewClient(t).Auth(nested.auth)
			authClient.POST("/projects").
				WithBody(map[string]interface{}{
					"project_name": "Test Project",
					"metadata":     map[string]string{"env": "test"},
				}).
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.AuthNotClient).
				HasMessage("only clients can access this endpoint")
		})

		t.Run("ListProjects", func(t *testing.T) {
			authClient := suite.NewClient(t).Auth(nested.auth)
			authClient.GET("/projects").
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.AuthNotClient).
				HasMessage("only clients can access this endpoint")
		})

		t.Run("GetProject", func(t *testing.T) {
			authClient := suite.NewClient(t).Auth(nested.auth)
			authClient.GET("/projects/" + user.ProjectID).
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.AuthNotClient).
				HasMessage("only clients can access this endpoint")
		})

		t.Run("UpdateProject", func(t *testing.T) {
			authClient := suite.NewClient(t).Auth(nested.auth)
			authClient.PATCH("/projects/" + user.ProjectID).
				WithBody(map[string]interface{}{
					"project_name": "Updated Project",
					"metadata":     map[string]string{"env": "prod"},
				}).
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.AuthNotClient).
				HasMessage("only clients can access this endpoint")
		})

		t.Run("GetProjectJWKS", func(t *testing.T) {
			jwksClient := suite.NewClient(t)
			jwksClient.GET("/projects/" + user.ProjectID + "/.well-known/jwks.json").
				Expect(http.StatusOK).
				JSON().Value("keys").
				Array().NotEmpty()
		})

		t.Run("DeleteProject", func(t *testing.T) {
			authClient := suite.NewClient(t).Auth(nested.auth)
			authClient.DELETE("/projects/" + user.ProjectID).
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.AuthNotClient).
				HasMessage("only clients can access this endpoint")
		})
	})
}
