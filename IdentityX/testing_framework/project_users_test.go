package testing

import (
	"GoAuth/internal/apierr"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func testProjectUsers(t *testing.T, suite *TestSuite) {
	// Create user in parent test context
	client := suite.NewClient(t)
	user := client.WithCredentials("client@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("test project")

	// Create with the same email as client to prove it's in a different environment
	projectUser := client.WithCredentials("client@mail.com", ValidPassword).
		ProjectRegister(user.projectID).
		ProjectLogin(user.projectID)

	t.Run("ProjectUsersRegister", func(t *testing.T) {
		var wrongFormatID = "wrong-format"
		t.Run("WrongFormatIDProjectRegister", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + wrongFormatID + "/register").
				WithBody(map[string]interface{}{
					"email":    user.email,
					"password": user.password,
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.RequestValidationError).
				HasMessage("Validation failed")
		})

		t.Run("InvalidProjectRegister", func(t *testing.T) {
			nonexistentID, err := uuid.NewV7()
			if err != nil {
				t.Fatal(err)
			}
			client := suite.NewClient(t)
			client.POST("/projects/" + nonexistentID.String() + "/register").
				WithBody(map[string]interface{}{
					"email":    user.email,
					"password": user.password,
				}).
				Expect(http.StatusBadRequest).
				HasMessage("invalid reference")
		})

		t.Run("ValidationProjectRegister", func(t *testing.T) {
			for _, spec := range ValidationTests {
				spec := spec
				t.Run(spec.Name, func(t *testing.T) {
					client := suite.NewClient(t)
					client.POST("/projects/" + user.projectID + "/register").
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

		t.Run("WeakPasswordValidationProjectRegister", func(t *testing.T) {
			for i, spec := range WeakPasswordTests {
				spec := spec
				i := i
				t.Run(spec.Name, func(t *testing.T) {
					client := suite.NewClient(t)
					client.POST("/projects/" + user.projectID + "/register").
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
			client.POST("/projects/" + user.projectID + "/register").
				WithBody(map[string]interface{}{
					"email":    user.email,
					"password": user.password,
				}).
				Expect(http.StatusConflict).
				HasErrID(apierr.AuthEmailAlreadyUsed).
				HasMessage("error registering user").
				TraceContains("email already in use")
		})

		t.Run("InvalidSchemaTypeRegister", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/"+user.projectID+"/register").
				WithQuery("schema_type", "invalid").
				WithBody(map[string]interface{}{
					"email":    "invalid_schema@email.com",
					"password": user.password,
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.SchemaInvalidSchemaType).
				HasMessage("invalid schema type")
		})

		t.Run("FlowIDSameAsTypeRegister", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/"+user.projectID+"/register").
				WithQuery("schema_type", "context").
				WithQuery("flow_id", "context").
				WithBody(map[string]interface{}{
					"email":    "flow_same_as_type@email.com",
					"password": user.password,
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.SchemaInvalidFlowID).
				HasMessage("flow id can't be the same as a schema type")
		})

		t.Run("MetadataRegisterOnCoreDenied", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + user.projectID + "/register").
				WithBody(map[string]interface{}{
					"email":    "metadata_denied@email.com",
					"password": user.password,
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
			client.WithCredentials("new@mail.com", ValidPassword).ProjectRegister(user.projectID)
		})
	})

	t.Run("ProjectUsersLogin", func(t *testing.T) {
		t.Run("WrongPassword", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + user.projectID + "/login").
				WithBody(map[string]string{
					"email":    projectUser.email,
					"password": "WrongPass123!",
				}).
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.AuthInvalidCredentials).
				HasMessage("invalid email or password")
		})

		t.Run("WrongEmail", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + user.projectID + "/login").
				WithBody(map[string]string{
					"email":    "wrong@mail.com",
					"password": projectUser.password,
				}).
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.AuthInvalidCredentials).
				HasMessage("invalid email or password")
		})

		t.Run("Success", func(t *testing.T) {
			client := suite.NewClient(t)
			client.WithCredentials(projectUser.email, projectUser.password).ProjectLogin(user.projectID)
		})

		t.Run("Logout", func(t *testing.T) {
			client := suite.NewClient(t)
			loggedInUser := client.WithCredentials(projectUser.email, projectUser.password).ProjectLogin(user.projectID)
			loggedInUser.Logout()

			// Try using revoked session
			loggedInUser.POST("/auth/logout").
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.SessionUnauthorized).
				HasMessage("session not found or revoked")
		})
	})

	sessionUser := client.WithCredentials("sessions@mail.com", ValidPassword).
		ProjectRegister(user.projectID).
		ProjectLogin(user.projectID)

	t.Run("ProjectUsersSession", func(t *testing.T) {
		t.Run("ListSessions", func(t *testing.T) {
			// Create a new client with subtest's t for the authenticated client
			// We need to preserve the auth context but use the new t
			authClient := suite.NewClient(t).WithAuth(sessionUser.auth)
			authClient.GET("/sessions").
				Expect(http.StatusOK).
				RequireDataArray().Length().IsEqual(1)
		})

		t.Run("MultipleLoginsSessions", func(t *testing.T) {
			// Create 3 more sessions (we already have 1)
			client2 := suite.NewClient(t)
			client2.WithCredentials(sessionUser.email, sessionUser.password).ProjectLogin(user.projectID)

			client3 := suite.NewClient(t)
			client3.WithCredentials(sessionUser.email, sessionUser.password).ProjectLogin(user.projectID)

			client4 := suite.NewClient(t)
			user4 := client4.WithCredentials(sessionUser.email, sessionUser.password).ProjectLogin(user.projectID)

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
			revokeOthers := client.WithCredentials("revoke-others-project@mail.com", ValidPassword).
				ProjectRegister(user.projectID).
				ProjectLogin(user.projectID)

			// Create more sessions
			suite.NewClient(t).WithCredentials(revokeOthers.email, revokeOthers.password).ProjectLogin(user.projectID)
			suite.NewClient(t).WithCredentials(revokeOthers.email, revokeOthers.password).ProjectLogin(user.projectID)

			revokeOthers.DELETE("/sessions/others").
				Expect(http.StatusOK).
				HasMessage("revoked sessions")

			// Should only have current session
			revokeOthers.GET("/sessions").
				Expect(http.StatusOK).
				RequireDataArray().Length().IsEqual(1)
		})

		t.Run("SessionInfo", func(t *testing.T) {
			client := suite.NewClient(t)
			infoUser := client.WithCredentials("session-me@mail.com", ValidPassword).
				ProjectRegister(user.projectID).
				ProjectLogin(user.projectID)

			data := infoUser.GET("/sessions/me").
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
						"project_id": AsString{user.projectID, AnyUUID{}},
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
			revoked := client.WithCredentials("revoke-all@mail.com", ValidPassword).
				ProjectRegister(user.projectID).
				ProjectLogin(user.projectID)

			// Create more sessions
			suite.NewClient(t).WithCredentials(revoked.email, revoked.password).ProjectLogin(user.projectID)
			suite.NewClient(t).WithCredentials(revoked.email, revoked.password).ProjectLogin(user.projectID)

			revoked.DELETE("/sessions").
				Expect(http.StatusOK).
				HasMessage("revoked sessions")

			// Session should be invalid
			revoked.GET("/sessions").
				Expect(http.StatusUnauthorized).
				HasErrID(apierr.SessionUnauthorized).
				HasMessage("session not found or revoked")
		})
	})

	t.Run("ProjectUserRefresh", func(t *testing.T) {
		client := suite.NewClient(t)
		refreshUser := client.WithCredentials("refresh@mail.com", ValidPassword).
			ProjectRegister(user.projectID).
			ProjectLogin(user.projectID)

		oldAccess := refreshUser.auth.AccessToken
		oldRefresh := refreshUser.auth.RefreshToken

		// Old tokens should be invalid
		oldClient := client.WithAuth(&AuthContext{
			AccessToken:  oldAccess,
			RefreshToken: oldRefresh,
		})

		refreshUser = refreshUser.Refresh()

		require.NotEqual(t, oldClient.auth.AccessToken, refreshUser.auth.AccessToken, "Access token should change after refresh")
		require.NotEqual(t, oldClient.auth.RefreshToken, refreshUser.auth.RefreshToken, "Refresh token should change after refresh")

		oldClient.GET("/sessions").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SessionUnauthorized).
			HasMessage("session not found or revoked")
	})

	t.Run("ProjectUsersProjects", func(t *testing.T) {
		client := suite.NewClient(t)
		nested := client.WithCredentials("nested_creator@mail.com", ValidPassword).
			ProjectRegister(user.projectID).
			ProjectLogin(user.projectID)

		t.Run("CreateProject", func(t *testing.T) {
			authClient := suite.NewClient(t).WithAuth(nested.auth)
			authClient.POST("/projects").
				WithBody(map[string]interface{}{
					"project_name": "Test Project",
					"metadata":     map[string]string{"env": "test"},
				}).
				Expect(http.StatusForbidden).
				HasErrID(apierr.AuthNotClient).
				HasMessage("only clients can access this endpoint")
		})

		t.Run("ListProjects", func(t *testing.T) {
			authClient := suite.NewClient(t).WithAuth(nested.auth)
			authClient.GET("/projects").
				Expect(http.StatusForbidden).
				HasErrID(apierr.AuthNotClient).
				HasMessage("only clients can access this endpoint")
		})

		t.Run("GetProject", func(t *testing.T) {
			authClient := suite.NewClient(t).WithAuth(nested.auth)
			authClient.GET("/projects/" + user.projectID).
				Expect(http.StatusForbidden).
				HasErrID(apierr.AuthNotClient).
				HasMessage("only clients can access this endpoint")
		})

		t.Run("UpdateProject", func(t *testing.T) {
			authClient := suite.NewClient(t).WithAuth(nested.auth)
			authClient.PATCH("/projects/" + user.projectID).
				WithBody(map[string]interface{}{
					"project_name": "Updated Project",
					"metadata":     map[string]string{"env": "prod"},
				}).
				Expect(http.StatusForbidden).
				HasErrID(apierr.AuthNotClient).
				HasMessage("only clients can access this endpoint")
		})

		t.Run("GetProjectJWKS", func(t *testing.T) {
			jwksClient := suite.NewClient(t)
			jwksClient.GET("/projects/" + user.projectID + "/.well-known/jwks.json").
				Expect(http.StatusOK).
				JSONObj().Value("keys").
				Array().NotEmpty()
		})

		t.Run("DeleteProject", func(t *testing.T) {
			authClient := suite.NewClient(t).WithAuth(nested.auth)
			authClient.DELETE("/projects/" + user.projectID).
				Expect(http.StatusForbidden).
				HasErrID(apierr.AuthNotClient).
				HasMessage("only clients can access this endpoint")
		})
	})

	t.Run("CryptographicIsolation", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("crypto@mail.com", ValidPassword).
			Register().
			Login().
			CreateProject("Crypto Project")

		projectID := user.projectID
		projectUser := client.WithCredentials("user@crypto.com", ValidPassword).
			ProjectRegister(projectID).
			ProjectLogin(projectID)

		accessToken := projectUser.auth.AccessToken

		// 1. Get Project JWKS (Raw)
		jwksResp := client.GET("/projects/" + projectID + "/.well-known/jwks.json").
			Expect(http.StatusOK).
			JSONObj()

		// Extract public key from JWKS
		xBase64 := jwksResp.Value("keys").Array().Value(0).Object().Value("x").String().Raw()
		xBytes, err := base64.RawURLEncoding.DecodeString(xBase64)
		require.NoError(t, err)
		projectPubKey := ed25519.PublicKey(xBytes)

		// 2. Get Global JWKS (Master - Wrapped in Data)
		globalJwksResp := client.GET("/.well-known/jwks.json").
			Expect(http.StatusOK).
			RequireDataObject()

		gxBase64 := globalJwksResp.Value("keys").Array().Value(0).Object().Value("x").String().Raw()
		gxBytes, err := base64.RawURLEncoding.DecodeString(gxBase64)
		require.NoError(t, err)
		masterPubKey := ed25519.PublicKey(gxBytes)

		// 3. Try verifying with Project Pub Key
		_, err = jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			return projectPubKey, nil
		})
		require.NoError(t, err, "Token for project %s should be verifiable with its own project JWKS", projectID)

		// 4. Try verifying with Master Pub Key
		_, err = jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			return masterPubKey, nil
		})
		require.Error(t, err, "Token MUST NOT be verifiable by Master Pub Key")
	})
}
