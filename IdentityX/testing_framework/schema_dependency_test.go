package testing

import (
	"net/http"
	"testing"
)

func testSchemaDependencies(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("dep_test@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("dep testing")

	projectID := user.projectID

	var schemaID string
	t.Run("CreateSchema", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "dependency-test-flow",
				"flow_id":     "dep-test",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

		schemaID = data.Object().Value("id").String().Raw()
	})

	t.Run("DraftVersion1", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusCreated)
	})

	t.Run("CreateBaseField", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/v1").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":      "base_field",
						"type":     "bool",
						"owner":    "user",
						"title":    "Base Field",
						"required": true,
						"mutable":  true,
						"position": 0,
					},
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("created fields")
	})

	t.Run("CreateDependentFieldInSeparateBatch", func(t *testing.T) {
		// This test expects failure with the BUG, and success with the FIX.
		// The bug is that 'base_field' won't be found because it's not in this batch.

		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/v1").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":      "dependent_field",
						"type":     "string",
						"owner":    "user",
						"title":    "Dependent Field",
						"required": false,
						"mutable":  true,
						"position": 1,
						"visibility_rules": []interface{}{
							map[string]interface{}{
								"depends_on_field_key": "base_field",
								"operator":             "equals",
								"value":                true,
							},
						},
					},
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("created fields")
	})
}
