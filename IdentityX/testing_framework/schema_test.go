package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func testSchemas(t *testing.T, suite *TestSuite) {
	client := suite.Client(t)
	user := client.User("schemas@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("schema testing")

	var projectID string
	projectID = user.ProjectID

	rid, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("Couldn't generate random uuid for test: %v", err)
	}

	t.Run("PublishSchemaRandomID", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + rid.String() + "/publish").
			Expect(http.StatusUnauthorized).
			MessageContains("cannot publish a schema you don't own").
			ExpectErrorID(apierr.SchemaNotOwnedByPrincipal)
	})

	var schemaID string
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
			MessageContains("cannot publish a schema version draft that doesn't exist").
			ExpectErrorID(apierr.SchemaVersionDraftDoesntExist)
	})

	var schemaVersion1ID string
	t.Run("DraftVersion", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusCreated).
			Data()

		data.Value("id").String().NotEmpty()
		data.Value("schema_id").String().IsEqual(schemaID)
		data.Value("version_number").IsNumber().IsEqual(1)

		schemaVersion1ID = data.Value("id").String().Raw()
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
		data.Value("current_version_id").IsEqual(schemaVersion1ID)
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
			Expect(http.StatusUnauthorized).
			MessageContains("new versions can only be drafted from published versions").
			ExpectErrorID(apierr.SchemaVersionDraftOnNonPublished)
	})

	t.Run("PublishVersionFieldsError", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusUnauthorized).
			MessageContains("cannot publish a schema version with no fields").
			ExpectErrorID(apierr.SchemaVersionPublishWithNoFields)
	})

	t.Run("CreateFieldsSamePosition", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/v1").
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
						"position":    0,
					},
				},
			}).
			Expect(http.StatusConflict).
			MessageContains("two fields can't occupy the same position").
			ExpectErrorID(apierr.FieldSamePositionForMultipleFields)
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
		id1 := data.Value(0).Object().Value("id").String().NotEmpty().Raw()
		data.Value(1).Object().Value("object_id").NotNull()
		id2 := data.Value(1).Object().Value("id").String().NotEmpty().Raw()

		if _, err := uuid.Parse(id1); err != nil {
			t.Fatalf("couldn't parse id from field matricula: %v", err)
		}
		if _, err := uuid.Parse(id2); err != nil {
			t.Fatalf("couldn't parse id from field curso: %v", err)
		}
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

	var schemaVersion2ID string
	t.Run("DraftVersion2", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusCreated).
			Data()

		data.Value("id").String().NotEmpty()
		data.Value("schema_id").String().IsEqual(schemaID)
		data.Value("version_number").IsNumber().IsEqual(2)

		schemaVersion2ID = data.Value("id").String().Raw()
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
		data.Value("status").String().IsEqual("published")
		data.Value("current_version_id").IsEqual(schemaVersion2ID)
	})

	t.Run("PublishVersion2NoChanges", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusBadRequest).
			MessageContains("cannot publish a version with no changes").
			ExpectErrorID(apierr.SchemaVersionNoChanges)
	})

	t.Run("AddFieldToV2FailKeyCheck", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/v2").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":         "período",
						"type":        "int",
						"owner":       "user",
						"title":       "Período Atual",
						"description": "O período da sua matéria mais avançada da grade",
						"required":    true,
						"mutable":     true,
						"position":    2,
					},
				},
			}).
			Expect(http.StatusBadRequest).
			MessageContains("field key must start with a lowercase letter and contain only lowercase letters, numbers, or underscores").
			ExpectErrorID(apierr.FieldInvalidCharactersInKey)
	})

	t.Run("AddFieldToV2Success", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/v2").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":         "periodo",
						"type":        "int",
						"owner":       "user",
						"title":       "Período Atual",
						"description": "O período da sua matéria mais avançada da grade",
						"required":    true,
						"mutable":     true,
						"position":    2,
					},
				},
			}).
			Expect(http.StatusCreated).
			MessageContains("created fields").
			DataArray()

		data.Length().IsEqual(1)
		data.Value(0).Object().Value("object_id").NotNull()
		id1 := data.Value(0).Object().Value("id").String().NotEmpty().Raw()

		if _, err := uuid.Parse(id1); err != nil {
			t.Fatalf("couldn't parse id from field periodo: %v", err)
		}
	})

	t.Run("PublishVersion2Success", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusOK).
			MessageContains("published schema version")
	})

	t.Run("GetSchemaVerbose", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		schema := authClient.GET("/projects/" + projectID + "/schemas/" + schemaID + "/verbose").
			Expect(http.StatusOK).
			Data()

		t.Log(schema)

		schema.Value("id").String().IsEqual(schemaID)
		schema.Value("project_id").String().IsEqual(projectID)
		schema.Value("title").String().IsEqual("scti-register-flow")
		schema.Value("flow_id").String().IsEqual("scti-register")
		schema.Value("type").String().IsEqual("context")
		schema.Value("status").String().IsEqual("published")
		schema.Value("current_version_id").String().IsEqual(schemaVersion2ID)
		schema.Value("created_at").String().NotEmpty()
		schema.Value("updated_at").String().NotEmpty()

		versions := schema.Value("versions").Array()
		versions.Length().IsEqual(2)

		version2 := versions.Value(0).Object()
		version1 := versions.Value(1).Object()

		fieldsV1 := version1.Value("fields").Array()
		fieldsV2 := version2.Value("fields").Array()

		fieldsV1.Length().IsEqual(2)
		fieldsV2.Length().IsEqual(3)

		fieldsV1field1ID := fieldsV1.Value(0).Object().Value("id").String().Raw()
		fieldsV1field2ID := fieldsV1.Value(1).Object().Value("id").String().Raw()

		fieldsV2field1ID := fieldsV2.Value(0).Object().Value("id").String().Raw()
		fieldsV2field2ID := fieldsV2.Value(1).Object().Value("id").String().Raw()

		if fieldsV1field1ID != fieldsV2field1ID {
			t.Fatalf("field id doesn't match between versions")
		}

		if fieldsV1field2ID != fieldsV2field2ID {
			t.Fatalf("field id doesn't match between versions")
		}
	})
}
