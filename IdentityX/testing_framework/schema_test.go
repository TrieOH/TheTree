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

	t.Run("PublishSchemaNoVersion", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusUnauthorized).
			MessageContains("cannot publish a schema with no versions").
			ExpectErrorID(apierr.SchemaNoPublishedVersion)
	})

	t.Run("PublishVersionNoDraft", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusUnauthorized).
			MessageContains("cannot publish a schema version draft that doesn't exists").
			ExpectErrorID(apierr.SchemaVersionDraftDoesntExist)
	})

	var schemaVersionID string
	t.Run("DraftVersion", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
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

	t.Run("PublishSchemaOnlyDraft", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusUnauthorized).
			MessageContains("cannot publish a schema with only draft versions").
			ExpectErrorID(apierr.SchemaHasOnlyDraftVersion)
	})

	t.Run("DraftVersionError", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusConflict).
			MessageContains("a draft schema version already exists").
			ExpectErrorID(apierr.SchemaVersionDraftAlreadyExists)
	})

	t.Run("PublishVersionFieldsError", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusUnauthorized).
			MessageContains("cannot publish a schema version with no fields").
			ExpectErrorID(apierr.SchemaVersionPublishWithNoFields)
	})

	t.Run("CreateFields", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/v1").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":         "matricula",
						"type":        "string",
						"owner":       "user",
						"title":       "Numero da Matrícula",
						"description": "Sua matrícula da UENF como aparece no sistema acadêmico",
						"placeholder": "20223200045",
						"required":    true,
						"mutable":     true,
						"position":    0,
					},
					map[string]interface{}{
						"key":         "curso",
						"type":        "string",
						"owner":       "user",
						"title":       "Curso de Matrícula",
						"description": "O curso que você está matrículado na UENF",
						"placeholder": "Ciência da Computação",
						"required":    true,
						"mutable":     true,
						"position":    1,
					},
				},
			}).
			Expect(http.StatusCreated).
			MessageContains("created fields").
			DataArray()

		data.Length().IsEqual(2)
		data.Value(0).Object().Value("object_id").NotNull()
		data.Value(0).Object().Value("id").NotNull()
		data.Value(1).Object().Value("object_id").NotNull()
		data.Value(1).Object().Value("id").NotNull()
	})

	t.Run("PublishVersionSuccess", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusOK).
			MessageContains("published schema version")
	})

	t.Run("PublishVersionAlreadyPublished", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusUnauthorized).
			MessageContains("cannot publish a schema version that isn't a draft").
			ExpectErrorID(apierr.SchemaVersionTryingToPublishPublished)
	})

	t.Run("PublishSchemaSuccess", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusOK).
			MessageContains("published schema")
	})

	t.Run("PublishSchemaAlreadyPublished", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusUnauthorized).
			MessageContains("cannot publish a schema that isn't a draft").
			ExpectErrorID(apierr.SchemaTryingToPublishPublished)
	})
}
