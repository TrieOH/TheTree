package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"
)

func testProjects(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("projects@mail.com", ValidPassword).Register().Login()

	var projectID string
	t.Run("CreateProject", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.POST("/projects").
			WithBody(map[string]interface{}{
				"project_name": "Test Project",
				"metadata":     map[string]string{"env": "test"},
			}).
			Expect(http.StatusCreated).
			HasMessage("Created project").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":           StoreString{Into: &projectID, Matcher: AnyUUID{}},
			"owner_id":     AnyUUID{},
			"project_name": "Test Project",
			"is_active":    true,
			"metadata": map[string]interface{}{
				"env": "test",
			},
			"created_at": AnyDate{},
			"updated_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("ValidationCreateProject", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects").
			WithBody(map[string]interface{}{
				"project_name": "", // Empty name
				"metadata":     map[string]string{"env": "test"},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.RequestValidationError).
			ValidationError("project_name is required")
	})

	t.Run("ListProjects", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.GET("/projects").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":           AsString{projectID, AnyUUID{}},
				"owner_id":     AnyUUID{},
				"project_name": "Test Project",
				"is_active":    true,
				"metadata": map[string]interface{}{
					"env": "test",
				},
				"created_at": AnyDate{},
				"updated_at": AnyDate{},
			},
		}

		Validate(t, data, spec)
	})

	t.Run("GetProject", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.GET("/projects/" + projectID).
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":           AsString{projectID, AnyUUID{}},
			"owner_id":     AnyUUID{},
			"project_name": "Test Project",
			"is_active":    true,
			"metadata": map[string]interface{}{
				"env": "test",
			},
			"created_at": AnyDate{},
			"updated_at": AnyDate{},
		}

		Validate(t, data, spec)
	})

	t.Run("CrossUserAccess", func(t *testing.T) {
		// Create a second user
		attacker := client.WithCredentials("attacker@mail.com", ValidPassword).Register().Login()
		attackerClient := suite.NewClient(t).WithAuth(attacker.auth)

		// Try to GET project owned by first user
		attackerClient.GET("/projects/" + projectID).
			Expect(http.StatusNotFound).
			HasErrID(apierr.DBNotFound).
			HasMessage("resource not found")

		// Try to UPDATE
		attackerClient.PATCH("/projects/" + projectID).
			WithBody(map[string]interface{}{
				"project_name": "Hacked",
			}).
			Expect(http.StatusNotFound).
			HasErrID(apierr.DBNotFound).
			HasMessage("resource not found")

		// Try to DELETE
		attackerClient.DELETE("/projects/" + projectID).
			Expect(http.StatusNotFound).
			HasErrID(apierr.ProjectNotFound).
			HasMessage("project not found")

		// Ensure it was NOT actually deleted from the perspective of the owner
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.GET("/projects/" + projectID).
			Expect(http.StatusOK)
	})

	t.Run("UpdateProject", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.PATCH("/projects/" + projectID).
			WithBody(map[string]interface{}{
				"project_name": "Updated Project",
				"metadata":     map[string]interface{}{"env": "prod"},
			}).
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":           AsString{projectID, AnyUUID{}},
			"owner_id":     AnyUUID{},
			"project_name": "Updated Project",
			"is_active":    true,
			"metadata": map[string]interface{}{
				"env": "prod",
			},
			"created_at": AnyDate{},
			"updated_at": AnyDate{},
		}

		Validate(t, data, spec)
	})

	t.Run("GetProjectJWKS", func(t *testing.T) {
		jwksClient := suite.NewClient(t)
		obj := jwksClient.GET("/projects/" + projectID + "/.well-known/jwks.json").
			Expect(http.StatusOK).
			JSONObj()

		obj.Value("keys").Array().NotEmpty()
	})

	t.Run("DeleteProject", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.DELETE("/projects/" + projectID).
			Expect(http.StatusOK).
			HasMessage("Deleted project")

		// Verify deletion
		authClient.GET("/projects").
			Expect(http.StatusOK).
			RequireDataArray().Length().IsEqual(0)
	})
}
