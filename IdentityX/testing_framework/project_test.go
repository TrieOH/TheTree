package testing

import (
	"net/http"
	"testing"
)

func testProjects(t *testing.T, suite *TestSuite) {
	client := suite.Client(t)
	user := client.User("projects@mail.com", ValidPassword).Register().Login()

	var projectID string
	t.Run("CreateProject", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		val := authClient.POST("/projects").
			WithBody(map[string]interface{}{
				"project_name": "Test Project",
				"metadata":     map[string]string{"env": "test"},
			}).
			Expect(http.StatusCreated).
			Value()

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
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects").
			WithBody(map[string]interface{}{
				"project_name": "", // Empty name
				"metadata":     map[string]string{"env": "test"},
			}).
			Expect(http.StatusBadRequest).
			ValidationError("(project_name) is required")
	})

	t.Run("ListProjects", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.GET("/projects").
			Expect(http.StatusOK).
			Value()

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
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.GET("/projects/" + projectID).
			Expect(http.StatusOK).
			Value()

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
		attacker := client.User("attacker@mail.com", ValidPassword).Register().Login()
		attackerClient := suite.Client(t).Auth(attacker.auth)

		// Try to GET project owned by first user
		attackerClient.GET("/projects/"+projectID).
			Expect(http.StatusNotFound).
			Error("go-auth-test", "resource not found")

		// Try to UPDATE
		attackerClient.PATCH("/projects/"+projectID).
			WithBody(map[string]interface{}{
				"project_name": "Hacked",
			}).
			Expect(http.StatusUnauthorized).
			Error("go-auth-test", "cannot update a project you don't own")

		// Try to DELETE
		attackerClient.DELETE("/projects/"+projectID).
			Expect(http.StatusUnauthorized).
			Success("go-auth-test", "cannot delete a project you don't own")

		// Ensure it was NOT actually deleted from the perspective of the owner
		authClient := suite.Client(t).Auth(user.auth)
		authClient.GET("/projects/" + projectID).
			Expect(http.StatusOK)
	})

	t.Run("UpdateProject", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.PATCH("/projects/" + projectID).
			WithBody(map[string]interface{}{
				"project_name": "Updated Project",
				"metadata":     map[string]string{"env": "prod"},
			}).
			Expect(http.StatusOK).
			Value()

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
