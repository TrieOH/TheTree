package testing

import (
	"GoAuth/internal/apierr"
	"fmt"
	"net/http"
	"testing"
)

func testProjectUsers(t *testing.T, suite *TestSuite) {
	// Create user in parent test context
	client := suite.Client(t)
	user := client.User("client@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("test project")

	// Create with the same email as client to prove it's in a different environment
	projectUser := client.User("client@mail.com", ValidPassword).
		ProjectRegister(user.ProjectID).
		ProjectLogin(user.ProjectID)

	t.Run("ProjectUsersRegister", func(t *testing.T) {
		var wrongFormatID = "wrong-format"
		t.Run("WrongFormatIDProjectRegister", func(t *testing.T) {
			client := suite.Client(t)
			client.POST("/projects/"+wrongFormatID+"/register").
				WithBody(map[string]interface{}{
					"email":    user.Email,
					"password": user.Password,
				}).
				Expect(http.StatusBadRequest).
				Error("go-auth-test", "invalid project id")
		})

		var nonexistentID = "d917a262-199a-41b1-af43-930ea8da1c75"
		t.Run("InvalidProjectRegister", func(t *testing.T) {
			client := suite.Client(t)
			client.POST("/projects/"+nonexistentID+"/register").
				WithBody(map[string]interface{}{
					"email":    user.Email,
					"password": user.Password,
				}).
				Expect(http.StatusBadRequest).
				Error("go-auth-test", "invalid reference").
				TraceContains("violates foreign key constraint")
		})

		t.Run("ValidationProjectRegister", func(t *testing.T) {
			for _, spec := range ValidationTests {
				spec := spec
				t.Run(spec.Name, func(t *testing.T) {
					client := suite.Client(t)
					client.POST("/projects/" + user.ProjectID + "/register").
						WithBody(map[string]interface{}{
							"email":    spec.Email,
							"password": spec.Pass,
						}).
						Expect(http.StatusBadRequest).
						ValidationError(spec.Errors...)
				})
			}
		})

		t.Run("ValidationProjectRegister", func(t *testing.T) {
			for i, spec := range WeakPasswordTests {
				spec := spec
				i := i
				t.Run(spec.Name, func(t *testing.T) {
					client := suite.Client(t)
					client.POST("/projects/" + user.ProjectID + "/register").
						WithBody(map[string]interface{}{
							"email":    fmt.Sprintf("weak%d@mail.com", i),
							"password": spec.Password,
						}).
						Expect(http.StatusBadRequest).
						ValidationError(spec.Errors...)
				})
			}
		})

		t.Run("DuplicateEmailProjectRegister", func(t *testing.T) {
			client := suite.Client(t)
			client.POST("/projects/"+user.ProjectID+"/register").
				WithBody(map[string]interface{}{
					"email":    user.Email,
					"password": user.Password,
				}).
				Expect(http.StatusConflict).
				Error("go-auth-test", "error registering user").
				TraceContains("email already in use")
		})

		t.Run("MetadataRegisterOnCoreDenied", func(t *testing.T) {
			client := suite.Client(t)
			client.POST("/projects/" + user.ProjectID + "/register").
				WithBody(map[string]interface{}{
					"email":    "metadata_denied@email.com",
					"password": user.Password,
					"custom_fields": map[string]interface{}{
						"curso": "Ciência da Computação",
					},
				}).
				Expect(http.StatusBadRequest).
				MessageContains("custom fields are not allowed for core schema").
				ExpectErrorID(apierr.SchemaMetadataNotAllowed)
		})

		t.Run("SuccessProjectRegister", func(t *testing.T) {
			client := suite.Client(t)
			client.User("new@mail.com", ValidPassword).ProjectRegister(user.ProjectID)
		})
	})

	t.Run("ProjectUsersLogin", func(t *testing.T) {
		t.Run("WrongPassword", func(t *testing.T) {
			client := suite.Client(t)
			client.POST("/projects/"+user.ProjectID+"/login").
				WithBody(map[string]string{
					"email":    projectUser.Email,
					"password": "WrongPass123!",
				}).
				Expect(http.StatusUnauthorized).
				Error("go-auth-test", "invalid email or password")
		})

		t.Run("WrongEmail", func(t *testing.T) {
			client := suite.Client(t)
			client.POST("/projects/"+user.ProjectID+"/login").
				WithBody(map[string]string{
					"email":    "wrong@mail.com",
					"password": projectUser.Password,
				}).
				Expect(http.StatusUnauthorized).
				Error("go-auth-test", "invalid email or password")
		})

		t.Run("Success", func(t *testing.T) {
			client := suite.Client(t)
			client.User(projectUser.Email, projectUser.Password).ProjectLogin(user.ProjectID)
		})

		t.Run("Logout", func(t *testing.T) {
			client := suite.Client(t)
			loggedInUser := client.User(projectUser.Email, projectUser.Password).ProjectLogin(user.ProjectID)
			loggedInUser.Logout()

			// Try using revoked session
			loggedInUser.AuthedClient().POST("/auth/logout").
				Expect(http.StatusUnauthorized).
				Error("AuthMW", "refresh token is revoked")
		})
	})

	sessionUser := client.User("sessions@mail.com", ValidPassword).
		ProjectRegister(user.ProjectID).
		ProjectLogin(user.ProjectID)

	t.Run("ProjectUsersSession", func(t *testing.T) {
		t.Run("ListSessions", func(t *testing.T) {
			// Create a new client with subtest's t for the authenticated client
			// We need to preserve the auth context but use the new t
			authClient := suite.Client(t).Auth(sessionUser.auth)
			arr := authClient.GET("/sessions").
				Expect(http.StatusOK).
				DataArray()

			arr.Length().IsEqual(1)
		})

		t.Run("MultipleLoginsSessions", func(t *testing.T) {
			// Create 3 more sessions (we already have 1)
			client2 := suite.Client(t)
			client2.User(sessionUser.Email, sessionUser.Password).ProjectLogin(user.ProjectID)

			client3 := suite.Client(t)
			client3.User(sessionUser.Email, sessionUser.Password).ProjectLogin(user.ProjectID)

			client4 := suite.Client(t)
			user4 := client4.User(sessionUser.Email, sessionUser.Password).ProjectLogin(user.ProjectID)

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
			revokeOthers := client.User("revoke-others-project@mail.com", ValidPassword).
				ProjectRegister(user.ProjectID).
				ProjectLogin(user.ProjectID)

			// Create more sessions
			suite.Client(t).User(revokeOthers.Email, revokeOthers.Password).ProjectLogin(user.ProjectID)
			suite.Client(t).User(revokeOthers.Email, revokeOthers.Password).ProjectLogin(user.ProjectID)

			revokeOthers.AuthedClient().DELETE("/sessions/others").
				Expect(http.StatusOK).
				Success("go-auth-test", "revoked sessions")

			// Should only have current session
			revokeOthers.AuthedClient().GET("/sessions").
				Expect(http.StatusOK).
				DataArray().
				Length().IsEqual(1)
		})

		t.Run("SessionInfo", func(t *testing.T) {
			client := suite.Client(t)
			infoUser := client.User("session-me@mail.com", ValidPassword).
				ProjectRegister(user.ProjectID).
				ProjectLogin(user.ProjectID)

			data := infoUser.AuthedClient().GET("/sessions/me").
				Expect(http.StatusOK).
				Data()

			data.Value("refresh_expire_date").IsNumber()

			access := data.Value("access").Object()
			access.Value("iss").String().IsEqual("GoAuth")
			access.Value("exp").IsNumber()
			access.Value("iat").IsNumber()
			access.Value("jti").String().NotEmpty()

			sub := access.Value("sub").Object()
			sub.Value("id").String().NotEmpty()
			sub.Value("email").String().IsEqual("session-me@mail.com")
			sub.Value("project_id").IsEqual(user.ProjectID)
			sub.Value("user_type").String().IsEqual("project")
			sub.Value("metadata").Object().IsEmpty()
			sub.Value("session_id").String().NotEmpty()
			sub.Value("user_agent").String().NotEmpty()
			sub.Value("user_ip").String().NotEmpty()
		})

		t.Run("RevokeAllSessions", func(t *testing.T) {
			client := suite.Client(t)
			revoked := client.User("revoke-all@mail.com", ValidPassword).
				ProjectRegister(user.ProjectID).
				ProjectLogin(user.ProjectID)

			// Create more sessions
			suite.Client(t).User(revoked.Email, revoked.Password).ProjectLogin(user.ProjectID)
			suite.Client(t).User(revoked.Email, revoked.Password).ProjectLogin(user.ProjectID)

			revoked.AuthedClient().DELETE("/sessions").
				Expect(http.StatusOK).
				Success("go-auth-test", "revoked sessions")

			// Session should be invalid
			revoked.AuthedClient().GET("/sessions").
				Expect(http.StatusUnauthorized).
				Error("AuthMW", "refresh token is revoked")
		})
	})

	t.Run("ProjectUserRefresh", func(t *testing.T) {
		client := suite.Client(t)
		refreshUSer := client.User("refresh@mail.com", ValidPassword).
			ProjectRegister(user.ProjectID).
			ProjectLogin(user.ProjectID)

		oldAccess := refreshUSer.auth.AccessToken
		oldRefresh := refreshUSer.auth.RefreshToken

		refreshUSer.Refresh()

		if oldAccess == refreshUSer.auth.AccessToken {
			t.Error("Access token should change after refresh")
		}

		if oldRefresh == refreshUSer.auth.RefreshToken {
			t.Error("Refresh token should change after refresh")
		}

		// Old tokens should be invalid
		oldClient := client.Auth(&AuthContext{
			AccessToken:  oldAccess,
			RefreshToken: oldRefresh,
		})

		oldClient.GET("/sessions").
			Expect(http.StatusUnauthorized)
	})

	t.Run("ProjectUsersProjects", func(t *testing.T) {
		client := suite.Client(t)
		nested := client.User("nested_creator@mail.com", ValidPassword).
			ProjectRegister(user.ProjectID).
			ProjectLogin(user.ProjectID)

		t.Run("CreateProject", func(t *testing.T) {
			authClient := suite.Client(t).Auth(nested.auth)
			authClient.POST("/projects").
				WithBody(map[string]interface{}{
					"project_name": "Test Project",
					"metadata":     map[string]string{"env": "test"},
				}).
				Expect(http.StatusUnauthorized).
				Error("ClientOnlyMW", "only clients can access this endpoint")
		})

		t.Run("ListProjects", func(t *testing.T) {
			authClient := suite.Client(t).Auth(nested.auth)
			authClient.GET("/projects").
				Expect(http.StatusUnauthorized).
				Error("ClientOnlyMW", "only clients can access this endpoint")
		})

		t.Run("GetProject", func(t *testing.T) {
			authClient := suite.Client(t).Auth(nested.auth)
			authClient.GET("/projects/"+user.ProjectID).
				Expect(http.StatusUnauthorized).
				Error("ClientOnlyMW", "only clients can access this endpoint")
		})

		t.Run("UpdateProject", func(t *testing.T) {
			authClient := suite.Client(t).Auth(nested.auth)
			authClient.PATCH("/projects/"+user.ProjectID).
				WithBody(map[string]interface{}{
					"project_name": "Updated Project",
					"metadata":     map[string]string{"env": "prod"},
				}).
				Expect(http.StatusUnauthorized).
				Error("ClientOnlyMW", "only clients can access this endpoint")
		})

		t.Run("GetProjectJWKS", func(t *testing.T) {
			jwksClient := suite.Client(t)
			obj := jwksClient.GET("/projects/" + user.ProjectID + "/.well-known/jwks.json").
				Expect(http.StatusOK).
				JSON()

			obj.Value("keys").Array().NotEmpty()
		})

		t.Run("DeleteProject", func(t *testing.T) {
			authClient := suite.Client(t).Auth(nested.auth)
			authClient.DELETE("/projects/"+user.ProjectID).
				Expect(http.StatusUnauthorized).
				Error("ClientOnlyMW", "only clients can access this endpoint")
		})
	})
}
