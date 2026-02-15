package testing

import (
	"GoAuth/internal/errx"
	"net/http"
	"testing"
)

func testApiKeys(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	owner := client.WithCredentials("api_owner@mail.com", ValidPassword).Register().Login()
	otherOwner := client.WithCredentials("api_other@mail.com", ValidPassword).Register().Login()

	// Owner creates two projects
	ownerA := owner.WithT(t).CreateProject("Project A")
	projectIDA := ownerA.ProjectID()
	ownerB := owner.WithT(t).CreateProject("Project B")
	projectIDB := ownerB.ProjectID()

	// Other owner creates a project
	otherC := otherOwner.WithT(t).CreateProject("Project C")
	projectIDC := otherC.ProjectID()

	var apiKeyA string
	t.Run("RotateApiKey_Success", func(t *testing.T) {
		authClient := ownerA.WithT(t)
		val := authClient.POST("/projects/" + projectIDA + "/api-keys/rotate").
			Expect(http.StatusOK).
			HasMessage("API key rotated").
			RequireDataValue()

		spec := map[string]interface{}{
			"api_key": StoreString{Into: &apiKeyA, Matcher: NotEmpty{}},
		}

		Validate(t, val, spec)
	})

	t.Run("RotateApiKey_Unauthorized", func(t *testing.T) {
		authClient := otherC.WithT(t)
		authClient.POST("/projects/" + projectIDA + "/api-keys/rotate").
			Expect(http.StatusNotFound).
			HasErrID(errx.ProjectNotFound)
	})

	t.Run("Rotation_InvalidatesOldKey", func(t *testing.T) {
		oldKey := apiKeyA
		authClient := ownerA.WithT(t)

		// Rotate again
		var newKey string
		val := authClient.POST("/projects/" + projectIDA + "/api-keys/rotate").
			Expect(http.StatusOK).
			RequireDataValue()

		Validate(t, val, map[string]interface{}{
			"api_key": StoreString{Into: &newKey, Matcher: NotEmpty{}},
		})

		// New key works
		suite.NewClient(t).WithApiKey(newKey).
			GET("/projects/" + projectIDA).
			Expect(http.StatusOK)

		// Old key fails
		suite.NewClient(t).WithApiKey(oldKey).
			GET("/projects/" + projectIDA).
			Expect(http.StatusUnauthorized).
			HasErrID(errx.AuthInvalidApiKey)

		apiKeyA = newKey // Update for subsequent tests
	})

	t.Run("CrossProjectIsolation_SameOwner", func(t *testing.T) {
		// Key for Project A should NOT be able to access Project B
		keyAClient := suite.NewClient(t).WithApiKey(apiKeyA)

		keyAClient.GET("/projects/" + projectIDB).
			Expect(http.StatusNotFound).
			HasErrID(errx.ProjectNotFound)
	})

	t.Run("CrossProjectIsolation_DifferentOwner", func(t *testing.T) {
		// Key for Project A should NOT be able to access Project C
		keyAClient := suite.NewClient(t).WithApiKey(apiKeyA)

		keyAClient.GET("/projects/" + projectIDC).
			Expect(http.StatusNotFound).
			HasErrID(errx.ProjectNotFound)
	})

	t.Run("OwnerCapabilities_ScopedToProject", func(t *testing.T) {
		keyAClient := suite.NewClient(t).WithApiKey(apiKeyA)

		// Should be able to create a schema in Project A
		keyAClient.POST("/projects/" + projectIDA + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"flow_id":     "api_test_flow",
				"title":       "API Test Flow",
			}).
			Expect(http.StatusCreated)

		// Should NOT be able to create a schema in Project B
		keyAClient.POST("/projects/" + projectIDB + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"flow_id":     "malicious_flow",
				"title":       "Malicious Flow",
			}).
			Expect(http.StatusNotFound) // Scoped to Project A, Project B is "not found"
	})

	t.Run("MalformedKeys", func(t *testing.T) {
		t.Run("WrongPrefix", func(t *testing.T) {
			suite.NewClient(t).WithApiKey("bad_prefix_uuid_secret").
				GET("/projects/" + projectIDA).
				Expect(http.StatusUnauthorized).
				HasErrID(errx.AuthInvalidApiKeyShape)
		})

		t.Run("WrongPartCount", func(t *testing.T) {
			suite.NewClient(t).WithApiKey("gk_too_few").
				GET("/projects/" + projectIDA).
				Expect(http.StatusUnauthorized).
				HasErrID(errx.AuthInvalidApiKeyShape)
		})

		t.Run("InvalidUUID", func(t *testing.T) {
			suite.NewClient(t).WithApiKey("gk_not-a-uuid_secret").
				GET("/projects/" + projectIDA).
				Expect(http.StatusUnauthorized).
				HasErrID(errx.AuthInvalidApiKeyShape)
		})
	})

	t.Run("RestrictedRoutes_NoApiKeys", func(t *testing.T) {
		keyAClient := suite.NewClient(t).WithApiKey(apiKeyA)

		keyAClient.GET("/sessions").
			Expect(http.StatusForbidden).
			HasErrID(errx.AuthApiKeyNotAllowed)

		keyAClient.POST("/auth/logout").
			Expect(http.StatusForbidden).
			HasErrID(errx.AuthApiKeyNotAllowed)
	})

	t.Run("FullCapabilities_ManageAndCheck", func(t *testing.T) {
		keyAClient := suite.NewClient(t).WithApiKey(apiKeyA).WithT(t)

		// 1. Register a Project User
		userEmail := "programmatic_user@mail.com"
		keyAClient.POST("/projects/" + projectIDA + "/register").
			WithBody(map[string]interface{}{
				"email":    userEmail,
				"password": ValidPassword,
			}).
			Expect(http.StatusCreated)

		// 2. Login as that user to get their ID
		puLogin := suite.NewClient(t).WithCredentials(userEmail, ValidPassword).
			ProjectLogin(projectIDA)
		userID := *puLogin.projectUserID

		// 3. Create Permission
		var permID string
		permVal := keyAClient.POST("/projects/" + projectIDA + "/permissions").
			WithBody(map[string]interface{}{
				"object": "resource:1",
				"action": "read",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

		Validate(t, permVal, map[string]interface{}{
			"id": StoreString{Into: &permID, Matcher: AnyUUID{}},
		})

		// 4. Create Role
		var roleID string
		roleVal := keyAClient.POST("/projects/" + projectIDA + "/roles").
			WithBody(map[string]interface{}{
				"name": "Manager",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

		Validate(t, roleVal, map[string]interface{}{
			"id": StoreString{Into: &roleID, Matcher: AnyUUID{}},
		})

		// 5. Add Permission to Role
		keyAClient.POST("/projects/" + projectIDA + "/roles/" + roleID + "/permissions/" + permID).
			Expect(http.StatusOK).
			HasMessage("Added permission to role")

		// 6. Assign Role to User
		keyAClient.POST("/projects/" + projectIDA + "/identities/" + userID + "/roles").
			WithBody(map[string]interface{}{
				"role_id": roleID,
			}).
			Expect(http.StatusOK).
			HasMessage("Added role to user")

		// 7. Perform Authz Check
		keyAClient.POST("/authz/check").
			WithBody(map[string]interface{}{
				"project_id": projectIDA,
				"entity_id":  userID,
				"object":     "resource:1",
				"action":     "read",
			}).
			Expect(http.StatusOK).
			HasMessage("Permission Granted")
	})

	t.Run("RevokeApiKey_Success", func(t *testing.T) {
		authClient := ownerA.WithT(t)
		authClient.DELETE("/projects/" + projectIDA + "/api-keys").
			Expect(http.StatusOK).
			HasMessage("API key revoked")

		// Verify it no longer works

		suite.NewClient(t).WithApiKey(apiKeyA).
			GET("/projects/" + projectIDA).
			Expect(http.StatusUnauthorized).
			HasErrID(errx.AuthInvalidApiKey)

	})

}
