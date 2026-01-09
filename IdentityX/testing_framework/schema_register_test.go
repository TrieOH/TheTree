package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"
)

func testSchemaRegister(t *testing.T, suite *TestSuite) {
	client := suite.Client(t)
	user := client.User("schemas_register@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("schema testing")

	projectID := user.ProjectID
	var schemaID string
	t.Run("Draft", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "scti",
				"flow_id":     "estudante",
			}).
			Expect(http.StatusCreated).
			Value()

		spec := map[string]interface{}{
			"id":                 StoreString{Into: &schemaID, Matcher: AnyUUID{}},
			"project_id":         AsString{projectID, AnyUUID{}},
			"title":              "scti",
			"flow_id":            "estudante",
			"type":               "context",
			"status":             "draft",
			"current_version_id": nil,
			"created_at":         AnyDate{},
			"updated_at":         AnyDate{},
		}

		Validate(t, data, spec)
	})

	var schemaVersion1ID string
	t.Run("DraftVersion", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusCreated).
			Value()

		spec := map[string]interface{}{
			"id":             StoreString{Into: &schemaVersion1ID, Matcher: AnyUUID{}},
			"schema_id":      AsString{schemaID, AnyUUID{}},
			"version_number": 1,
		}

		Validate(t, data, spec)
	})

	t.Run("CheckSchemaVersion", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.GET("/projects/" + projectID + "/schemas/" + schemaID).
			Expect(http.StatusOK).
			Value()

		spec := map[string]interface{}{
			"id":                 AsString{schemaID, AnyUUID{}},
			"project_id":         AsString{projectID, AnyUUID{}},
			"title":              "scti",
			"flow_id":            "estudante",
			"type":               "context",
			"status":             "draft",
			"current_version_id": AsString{schemaVersion1ID, AnyUUID{}},
		}

		Validate(t, data, spec)
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
			Value()

		spec := []interface{}{
			map[string]interface{}{
				"object_id": AnyUUID{},
				"id":        AnyUUID{},
			},
			map[string]interface{}{
				"object_id": AnyUUID{},
				"id":        AnyUUID{},
			},
			map[string]interface{}{
				"object_id": AnyUUID{},
				"id":        AnyUUID{},
			},
		}

		Validate(t, data, spec)
	})

	t.Run("PublishVersionSuccess", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusOK).
			MessageContains("published schema version")
	})

	t.Run("PublishSchemaSuccess", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusOK).
			MessageContains("published schema")
	})

	t.Run("RegisterOnSchemaNoCustomFields", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "client@email.com",
				"password": ValidPassword,
			}).
			Expect(http.StatusBadRequest).
			MessageContains("the schema custom fields are required on a schema register").
			ExpectErrorID(apierr.RequestMissingSchemaCustomFields)
	})

	t.Run("RegisterOnSchemaEmptyCustomFields", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":         "client@email.com",
				"password":      ValidPassword,
				"custom_fields": map[string]interface{}{},
			}).
			Expect(http.StatusBadRequest).
			MessageContains("missing required field").
			ExpectErrorID(apierr.FieldRequiredMissing)
	})

	t.Run("RegisterOnSchemaNoCursoField", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "client@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"matricula": "20221100033",
				},
			}).
			Expect(http.StatusBadRequest).
			MessageContains("missing required field").
			ExpectErrorID(apierr.FieldRequiredMissing)
	})

	t.Run("RegisterOnSchemaNoMatriculaField", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "client@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"curso": "Ciência da Computação",
				},
			}).
			Expect(http.StatusBadRequest).
			MessageContains("missing required field").
			ExpectErrorID(apierr.FieldRequiredMissing)
	})

	t.Run("RegisterOnSchemaUnknownField", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "client@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"valor": "4",
				},
			}).
			Expect(http.StatusBadRequest).
			MessageContains("unknown custom field").
			ExpectErrorID(apierr.FieldNotDefinedInSchema)
	})

	t.Run("RegisterOnSchemaWrongTypeStringOnInt", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "client@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"periodo": "abc",
				},
			}).
			Expect(http.StatusBadRequest).
			MessageContains("invalid field type").
			ExpectErrorID(apierr.FieldTypeMismatch)
	})

	t.Run("RegisterOnSchemaWrongTypeIntOnString", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "client@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"matricula": 20221100033,
				},
			}).
			Expect(http.StatusBadRequest).
			MessageContains("invalid field type").
			ExpectErrorID(apierr.FieldTypeMismatch)
	})

	t.Run("RegisterOnSchemaSuccess", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "client@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"matricula": "20221100033",
					"curso":     "Ciência da Computação",
					"periodo":   4,
				},
			}).
			Expect(http.StatusCreated).
			MessageContains("Registered user")
	})

	t.Run("SchemaUserSessionInfo", func(t *testing.T) {
		client := suite.Client(t)
		schemaUser := client.User("client@email.com", ValidPassword).ProjectLogin(user.ProjectID)
		data := schemaUser.AuthedClient().GET("/sessions/me").
			Expect(http.StatusOK).
			Value()

		spec := map[string]interface{}{
			"refresh_expire_date": AnyNumber{},
			"access": map[string]interface{}{
				"iss": "GoAuth",
				"exp": AnyNumber{},
				"iat": AnyNumber{},
				"jti": AnyUUID{},
				"sub": map[string]interface{}{
					"id":         AnyUUID{},
					"email":      "client@email.com",
					"project_id": projectID,
					"user_type":  "project",
					"session_id": AnyUUID{},
					"user_agent": AnyString{},
					"user_ip":    AnyString{},
					"metadata": map[string]interface{}{
						"context": map[string]interface{}{
							"estudante": map[string]interface{}{
								"schema_id":         schemaID,
								"schema_version_id": schemaVersion1ID,
								"curso":             "Ciência da Computação",
								"matricula":         "20221100033",
								"periodo":           4,
							},
						},
					},
				},
			},
		}

		Validate(t, data, spec)
	})
}
