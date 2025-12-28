package testing

import (
	"fmt"
	"net/http"
	"testing"
)

// ============================================================================
// TEST SPECS - Declarative test definitions
// ============================================================================

type ValidationSpec struct {
	Name   string
	Email  string
	Pass   string
	Errors []string
}

type PasswordSpec struct {
	Name     string
	Password string
	Errors   []string
}

// Test data
var (
	ValidPassword = "Str0ngP4$$!"

	ValidationTests = []ValidationSpec{
		{
			Name:   "NoEmail",
			Email:  "",
			Pass:   ValidPassword,
			Errors: []string{"(email) is required"},
		},
		{
			Name:   "InvalidEmail",
			Email:  "not-an-email",
			Pass:   ValidPassword,
			Errors: []string{"valid email address"},
		},
		{
			Name:   "NoPassword",
			Email:  "test@mail.com",
			Pass:   "",
			Errors: []string{"(password) is required"},
		},
	}

	WeakPasswordTests = []PasswordSpec{
		{"OnlyLetters", "abc", []string{"uppercase", "number", "symbol"}},
		{"LettersNumber", "abc3", []string{"uppercase", "symbol"}},
		{"LettersSymbol", "abc#", []string{"uppercase", "number"}},
		{"LettersUppercase", "Abc", []string{"number", "symbol"}},
		{"NoNumber", "Abc#", []string{"number"}},
		{"NoSymbol", "Abc3", []string{"symbol"}},
		{"TooShort", "Abc#3", []string{"at least 8 characters"}},
	}
)

// ============================================================================
// ACTUAL TESTS - Clean and readable
// ============================================================================

func TestGoAuth(t *testing.T) {
	suite := NewTestSuite(t)

	t.Run("Register", func(t *testing.T) {
		testRegister(t, suite)
	})

	t.Run("Login", func(t *testing.T) {
		testLogin(t, suite)
	})

	t.Run("Sessions", func(t *testing.T) {
		testSessions(t, suite)
	})

	t.Run("Refresh", func(t *testing.T) {
		testRefresh(t, suite)
	})

	t.Run("Projects", func(t *testing.T) {
		testProjects(t, suite)
	})
}

func testRegister(t *testing.T, suite *TestSuite) {
	t.Run("Validation", func(t *testing.T) {
		for _, spec := range ValidationTests {
			spec := spec // capture range variable
			t.Run(spec.Name, func(t *testing.T) {
				client := suite.Client(t)
				client.POST("/auth/register").
					WithBody(map[string]string{
						"email":    spec.Email,
						"password": spec.Pass,
					}).
					Expect(http.StatusBadRequest).
					ValidationError(spec.Errors...)
			})
		}
	})

	t.Run("WeakPasswords", func(t *testing.T) {
		for i, spec := range WeakPasswordTests {
			spec := spec // capture range variable
			i := i       // capture range variable
			t.Run(spec.Name, func(t *testing.T) {
				client := suite.Client(t)
				client.POST("/auth/register").
					WithBody(map[string]string{
						"email":    fmt.Sprintf("weak%d@mail.com", i),
						"password": spec.Password,
					}).
					Expect(http.StatusBadRequest).
					ValidationError(spec.Errors...)
			})
		}
	})

	t.Run("Success", func(t *testing.T) {
		client := suite.Client(t)
		client.User("new@mail.com", ValidPassword).Register()
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		client := suite.Client(t)
		email := "duplicate@mail.com"
		client.User(email, ValidPassword).Register()

		client.POST("/auth/register").
			WithBody(map[string]string{
				"email":    email,
				"password": ValidPassword,
			}).
			Expect(http.StatusConflict).
			Error("go-auth-test", "error registering user")
	})
}

func testLogin(t *testing.T, suite *TestSuite) {
	// Create user in parent test context
	client := suite.Client(t)
	user := client.User("login@mail.com", ValidPassword).Register()

	t.Run("WrongPassword", func(t *testing.T) {
		client := suite.Client(t)
		client.POST("/auth/login").
			WithBody(map[string]string{
				"email":    user.Email,
				"password": "WrongPass123!",
			}).
			Expect(http.StatusUnauthorized).
			Error("go-auth-test", "invalid email or password")
	})

	t.Run("WrongEmail", func(t *testing.T) {
		client := suite.Client(t)
		client.POST("/auth/login").
			WithBody(map[string]string{
				"email":    "wrong@mail.com",
				"password": user.Password,
			}).
			Expect(http.StatusUnauthorized).
			Error("go-auth-test", "invalid email or password")
	})

	t.Run("Success", func(t *testing.T) {
		client := suite.Client(t)
		client.User(user.Email, user.Password).Login()
	})

	t.Run("Logout", func(t *testing.T) {
		client := suite.Client(t)
		loggedInUser := client.User(user.Email, user.Password).Login()
		loggedInUser.Logout()

		// Try using revoked session
		loggedInUser.AuthedClient().POST("/auth/logout").
			Expect(http.StatusUnauthorized).
			Error("AuthMW", "refresh token is revoked")
	})
}

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

	t.Run("RevokeAllSessions", func(t *testing.T) {
		client := suite.Client(t)
		user := client.User("revoke-all@mail.com", ValidPassword).Register().Login()

		user.AuthedClient().DELETE("/sessions").
			Expect(http.StatusOK).
			Success("go-auth-test", "revoked sessions")

		// Session should be invalid
		user.AuthedClient().GET("/sessions").
			Expect(http.StatusUnauthorized).
			Error("AuthMW", "refresh token is revoked")
	})
}

func testRefresh(t *testing.T, suite *TestSuite) {
	client := suite.Client(t)
	user := client.User("refresh@mail.com", ValidPassword).Register().Login()

	oldAccess := user.auth.AccessToken
	oldRefresh := user.auth.RefreshToken

	user.Refresh()

	if oldAccess == user.auth.AccessToken {
		t.Error("Access token should change after refresh")
	}

	if oldRefresh == user.auth.RefreshToken {
		t.Error("Refresh token should change after refresh")
	}

	// Old tokens should be invalid
	oldClient := client.Auth(&AuthContext{
		AccessToken:  oldAccess,
		RefreshToken: oldRefresh,
	})

	oldClient.GET("/sessions").
		Expect(http.StatusUnauthorized)
}

func testProjects(t *testing.T, suite *TestSuite) {
	client := suite.Client(t)
	user := client.User("projects@mail.com", ValidPassword).Register().Login()

	var projectID string

	t.Run("CreateProject", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects").
			WithBody(map[string]interface{}{
				"project_name": "Test Project",
				"metadata":     map[string]string{"env": "test"},
			}).
			Expect(http.StatusCreated).
			Data()

		projectID = data.Value("id").String().Raw()
		data.Value("project_name").String().IsEqual("Test Project")
	})

	t.Run("ListProjects", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		arr := authClient.GET("/projects").
			Expect(http.StatusOK).
			DataArray()

		arr.Length().IsEqual(1)
		arr.Value(0).Object().Value("id").IsEqual(projectID)
	})

	t.Run("GetProject", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.GET("/projects/" + projectID).
			Expect(http.StatusOK).
			Data()

		data.Value("id").String().IsEqual(projectID)
		data.Value("project_name").String().IsEqual("Test Project")
	})

	t.Run("UpdateProject", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.PATCH("/projects/" + projectID).
			WithBody(map[string]interface{}{
				"project_name": "Updated Project",
				"metadata":     map[string]string{"env": "prod"},
			}).
			Expect(http.StatusOK).
			Data()

		data.Value("project_name").String().IsEqual("Updated Project")
	})

	t.Run("GetProjectKeys", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.GET("/projects/" + projectID + "/keys").
			Expect(http.StatusOK).
			Data()

		data.Value("pub_key").String().NotEmpty()
		data.Value("priv_key").NotNull()
	})

	t.Run("GetProjectJWKS", func(t *testing.T) {
		jwksClient := suite.Client(t)
		obj := jwksClient.GET("/projects/" + projectID + "/.well-known/jwks.json").
			Expect(http.StatusOK).
			JSON()

		obj.Value("keys").Array().NotEmpty()
	})

	t.Run("DeleteProject", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.DELETE("/projects/"+projectID).
			Expect(http.StatusOK).
			Success("go-auth-test", "Deleted project")

		// Verify deletion
		authClient.GET("/projects").
			Expect(http.StatusOK).
			DataArray().
			Length().IsEqual(0)
	})
}
