package testing

import (
	"GoAuth/internal/errx"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
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
			HasErrID(errx.RequestValidationError).
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

	t.Run("ListProjectUsers", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)

		// Register a user to this project first
		projectUserEmail := "projectuser@mail.com"
		suite.NewClient(t).POST("/projects/" + projectID + "/register").
			WithBody(map[string]interface{}{
				"email":    projectUserEmail,
				"password": ValidPassword,
			}).
			Expect(http.StatusCreated)

		data := authClient.GET("/projects/" + projectID + "/users").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AnyUUID{},
				"project_id": AsString{projectID, AnyUUID{}},
				"email":      projectUserEmail,
				"user_type":  "project",
				"is_active":  true,
				"created_at": AnyDate{},
				"updated_at": AnyDate{},
			},
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
			HasErrID(errx.SQLNotFound).
			HasMessage("project not found")

		// Try to UPDATE
		attackerClient.PATCH("/projects/" + projectID).
			WithBody(map[string]interface{}{
				"project_name": "Hacked",
			}).
			Expect(http.StatusNotFound).
			HasErrID(errx.SQLNotFound).
			HasMessage("project not found")

		// Try to DELETE
		attackerClient.DELETE("/projects/" + projectID).
			Expect(http.StatusNotFound).
			HasErrID(errx.ProjectNotFound).
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
		authClient := suite.NewClient(t).WithAuth(user.auth)
		obj := authClient.GET("/projects/" + projectID + "/.well-known/jwks.json").
			Expect(http.StatusOK).
			JSONObj()

		obj.Value("keys").Array().NotEmpty()

		// Verify decoding (Client gets it and decodes)
		xBase64 := obj.Value("keys").Array().Value(0).Object().Value("x").String().Raw()
		xBytes, err := base64.RawURLEncoding.DecodeString(xBase64)
		require.NoError(t, err)
		pubKey := parseJWKXToEd25519PublicKey(t, xBytes)
		require.NotNil(t, pubKey)

		// Unauthenticated access denial
		client.GET("/projects/" + projectID + "/.well-known/jwks.json").
			Expect(http.StatusUnauthorized)
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
