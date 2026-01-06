package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"
)

func testSchemas(t *testing.T, suite *TestSuite) {
	client := suite.Client(t)
	user := client.User("schemas@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("schema testing")

	var schemaID string
	var projectID string
	projectID = user.ProjectID
	t.Run("Draft", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "scti-register-flow",
				"flow_id":     "scti-register",
			}).
			Expect(http.StatusCreated).
			Data()

		data.Value("project_id").String().IsEqual(projectID)
		data.Value("title").String().IsEqual("scti-register-flow")
		data.Value("flow_id").String().IsEqual("scti-register")
		data.Value("id").String().NotEmpty()
		data.Value("type").String().IsEqual("context")
		data.Value("status").String().IsEqual("draft")
		data.Value("current_version_id").IsNull()

		schemaID = data.Value("id").String().Raw()
	})

	var schemaVersionID string
	t.Run("DraftVersion", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas/versions").
			WithBody(map[string]interface{}{
				"schema_id": schemaID,
			}).
			Expect(http.StatusCreated).
			Data()

		data.Value("id").String().NotEmpty()
		data.Value("schema_id").String().IsEqual(schemaID)
		data.Value("version_number").IsNumber().IsEqual(1)

		schemaVersionID = data.Value("id").String().Raw()
	})

	t.Run("CheckSchemaVersion", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.GET("/projects/" + projectID + "/schemas/" + schemaID).
			Expect(http.StatusOK).
			Data()

		data.Value("project_id").String().IsEqual(projectID)
		data.Value("title").String().IsEqual("scti-register-flow")
		data.Value("flow_id").String().IsEqual("scti-register")
		data.Value("id").String().NotEmpty()
		data.Value("type").String().IsEqual("context")
		data.Value("status").String().IsEqual("draft")
		data.Value("current_version_id").IsEqual(schemaVersionID)
	})

	t.Run("DraftVersionError", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/versions").
			WithBody(map[string]interface{}{
				"schema_id": schemaID,
			}).
			Expect(http.StatusConflict).
			MessageContains("a draft schema version already exists").
			ExpectErrorID(apierr.SchemaVersionDraftAlreadyExists)
	})
}
