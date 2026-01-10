package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"
)

func testSchemaRegister(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("schemas_register@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("schema testing")

	projectID := user.projectID
	var schemaID string
	t.Run("Draft", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "scti",
				"flow_id":     "estudante",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

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
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusCreated).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":             StoreString{Into: &schemaVersion1ID, Matcher: AnyUUID{}},
			"schema_id":      AsString{schemaID, AnyUUID{}},
			"version_number": 1,
		}

		Validate(t, data, spec)
	})

	t.Run("CheckSchemaVersion", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.GET("/projects/" + projectID + "/schemas/" + schemaID).
			Expect(http.StatusOK).
			RequireDataValue()

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
		authClient := suite.NewClient(t).WithAuth(user.auth)
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
					map[string]interface{}{
						"key":         "ativo",
						"type":        "bool",
						"owner":       "user",
						"title":       "Ativo",
						"description": "Se o aluno está ativo",
						"required":    false,
						"mutable":     true,
						"position":    3,
					},
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("created fields").
			RequireDataValue()

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
			map[string]interface{}{
				"object_id": AnyUUID{},
				"id":        AnyUUID{},
			},
		}

		Validate(t, data, spec)
	})

	t.Run("PublishVersionSuccess", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusOK).
			HasMessage("published schema version")
	})

	t.Run("PublishSchemaSuccess", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusOK).
			HasMessage("published schema")
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
			HasErrID(apierr.RequestMissingSchemaCustomFields).
			HasMessage("the schema custom fields are required on a schema register")
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
			HasErrID(apierr.FieldRequiredMissing).
			HasMessage("missing required field")
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
			HasErrID(apierr.FieldRequiredMissing).
			HasMessage("missing required field")
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
			HasErrID(apierr.FieldRequiredMissing).
			HasMessage("missing required field")
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
			HasErrID(apierr.FieldNotDefinedInSchema).
			HasMessage("unknown custom field")
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
			HasErrID(apierr.FieldTypeMismatch).
			HasMessage("invalid field type")
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
			HasErrID(apierr.FieldTypeMismatch).
			HasMessage("invalid field type")
	})

	t.Run("RegisterOnSchemaTypeFloatOnInt", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "float_on_int@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"matricula": "20221100033",
					"curso":     "Ciência da Computação",
					"periodo":   4.5,
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.FieldTypeMismatch).
			HasMessage("invalid field type")
	})

	t.Run("RegisterOnSchemaTypeStringOnBool", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "string_on_bool@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"matricula": "20221100033",
					"curso":     "Ciência da Computação",
					"periodo":   4,
					"ativo":     "true",
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.FieldTypeMismatch).
			HasMessage("invalid field type")
	})

	t.Run("RegisterOnSchemaTypeIntOnBool", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "int_on_bool@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"matricula": "20221100033",
					"curso":     "Ciência da Computação",
					"periodo":   4,
					"ativo":     1,
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.FieldTypeMismatch).
			HasMessage("invalid field type")
	})

	t.Run("RegisterOnSchemaTypeBoolOnString", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "bool_on_string@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"matricula": true,
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.FieldTypeMismatch).
			HasMessage("invalid field type")
	})

	t.Run("RegisterOnSchemaTypeFloatZeroOnInt", func(t *testing.T) {
		// Should succeed because 4.0 is a valid integer representation in JSON
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "float_zero@email.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"matricula": "20221100033",
					"curso":     "Ciência da Computação",
					"periodo":   4.0,
					"ativo":     true,
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
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
			HasMessage("Registered user")
	})

	t.Run("RegisterOnSchemaDuplicateEmail", func(t *testing.T) {
		client.POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":    "client@email.com", // Same email as above
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"matricula": "20221100033",
					"curso":     "Ciência da Computação",
					"periodo":   4,
				},
			}).
			Expect(http.StatusConflict).
			HasErrID(apierr.AuthEmailAlreadyUsed).
			HasMessage("error registering user").
			TraceContains("email already in use")
	})

	t.Run("SchemaUserSessionInfo", func(t *testing.T) {
		client := suite.NewClient(t)
		schemaUser := client.WithCredentials("client@email.com", ValidPassword).ProjectLogin(user.projectID)
		data := schemaUser.GET("/sessions/me").
			Expect(http.StatusOK).
			RequireDataValue()

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

	t.Run("SchemaStateEdgeCases", func(t *testing.T) {
		client := suite.NewClient(t)
		// New user/project for this isolation
		user := client.WithCredentials("schema_state@mail.com", ValidPassword).
			Register().
			Login().
			CreateProject("Schema State Project")

		projectID := user.projectID
		authClient := suite.NewClient(t).WithAuth(user.auth)

		// 1. Create a schema but don't create any version
		var schemaID string
		data := authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "No Version Schema",
				"flow_id":     "noversion",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

		schemaID = data.Path("$.id").String().Raw()

		t.Run("RegisterFailsNoVersion", func(t *testing.T) {
			client.POST("/projects/"+projectID+"/register").
				WithQuery("schema_type", "context").
				WithQuery("flow_id", "noversion").
				WithBody(map[string]interface{}{
					"email":         "user@noversion.com",
					"password":      ValidPassword,
					"custom_fields": map[string]interface{}{},
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.SchemaNoPublishedVersion).
				HasMessage("schema has no published version")
		})

		// 2. Create a version but don't publish it
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusCreated)

		// Add a field so it could potentially be published later
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/v1").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":      "test",
						"type":     "string",
						"owner":    "user",
						"title":    "Test",
						"position": 0,
						"required": true,
					},
				},
			}).
			Expect(http.StatusCreated)

		t.Run("RegisterFailsVersionNotPublished", func(t *testing.T) {
			// Even with a version draft, it's not the "current_version" yet because it's not published
			client.POST("/projects/"+projectID+"/register").
				WithQuery("schema_type", "context").
				WithQuery("flow_id", "noversion").
				WithBody(map[string]interface{}{
					"email":    "user@notpublished.com",
					"password": ValidPassword,
					"custom_fields": map[string]interface{}{
						"test": "val",
					},
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.ProjectUserRegisterOnSchemaDraft).
				HasMessage("can't register to a draft schema")
		})

		// 3. Publish the version but not the schema
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusOK)

		t.Run("RegisterFailsSchemaNotPublished", func(t *testing.T) {
			// Schema is still in 'draft' status
			client.POST("/projects/"+projectID+"/register").
				WithQuery("schema_type", "context").
				WithQuery("flow_id", "noversion").
				WithBody(map[string]interface{}{
					"email":    "user@schemadraft.com",
					"password": ValidPassword,
					"custom_fields": map[string]interface{}{
						"test": "val",
					},
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.ProjectUserRegisterOnSchemaDraft).
				HasMessage("can't register to a draft schema")
		})
	})

	t.Run("UnimplementedFieldTypes", func(t *testing.T) {
		clientU := suite.NewClient(t)
		userU := clientU.WithCredentials("unimplemented@mail.com", ValidPassword).
			Register().
			Login().
			CreateProject("Unimplemented Project")

		projectID := userU.projectID
		authClient := suite.NewClient(t).WithAuth(userU.auth)

		// Create Schema with 'email' type (valid ENUM but unimplemented in validator)
		var schemaIDU string
		data := authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "Email Schema",
				"flow_id":     "emailsync",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

		Validate(t, data, map[string]interface{}{
			"id": StoreString{Into: &schemaIDU, Matcher: AnyUUID{}},
		})

		authClient.POST("/projects/" + projectID + "/schemas/" + schemaIDU + "/versions/draft").
			Expect(http.StatusCreated)

		authClient.POST("/projects/" + projectID + "/schemas/" + schemaIDU + "/v1").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":      "contact",
						"type":     "email",
						"owner":    "user",
						"title":    "Contact Email",
						"position": 0,
						"required": true,
					},
				},
			}).
			Expect(http.StatusCreated)

		authClient.POST("/projects/" + projectID + "/schemas/" + schemaIDU + "/versions/publish").
			Expect(http.StatusOK)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaIDU + "/publish").
			Expect(http.StatusOK)

		t.Run("RegisterFailsDueToUnimplementedType", func(t *testing.T) {
			// This will fail because validateFieldType defaults to false for 'email'
			clientU.POST("/projects/"+projectID+"/register").
				WithQuery("schema_type", "context").
				WithQuery("flow_id", "emailsync").
				WithBody(map[string]interface{}{
					"email":    "user@unimplemented.com",
					"password": ValidPassword,
					"custom_fields": map[string]interface{}{
						"contact": "test@example.com",
					},
				}).
				Expect(http.StatusBadRequest).
				HasErrID(apierr.FieldTypeMismatch).
				HasMessage("invalid field type")
		})
	})
}
