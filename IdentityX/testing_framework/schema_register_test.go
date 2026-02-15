package testing

import (
	"GoAuth/internal/errx"
	"encoding/json"
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
			HasErrID(errx.RequestMissingSchemaCustomFields).
			HasMessage("schema custom fields are required on a schema register")
	})

	t.Run("RegisterOnSchemaEmptyCustomFields", func(t *testing.T) {
		client.WithT(t).POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "estudante").
			WithBody(map[string]interface{}{
				"email":         "client@email.com",
				"password":      ValidPassword,
				"custom_fields": map[string]interface{}{},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			HasMessage("error validating form for schema register").
			TraceContains("form missing required field: matricula", "form missing required field: curso", "form missing required field: periodo")
	})

	t.Run("RegisterOnSchemaNoCursoField", func(t *testing.T) {
		client.WithT(t).POST("/projects/"+projectID+"/register").
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
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			HasMessage("error validating form for schema register").
			TraceContains("form missing required field: curso", "form missing required field: periodo")
	})

	t.Run("RegisterOnSchemaNoMatriculaField", func(t *testing.T) {
		client.WithT(t).POST("/projects/"+projectID+"/register").
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
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			HasMessage("error validating form for schema register").
			TraceContains("form missing required field: matricula", "form missing required field: periodo")
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
					"bing":  "bong",
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			HasMessage("error validating form for schema register").
			TraceContains("form missing required field: matricula",
				"form missing required field: curso",
				"form missing required field: periodo")
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
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			HasMessage("error validating form for schema register").
			TraceContains(
				"form missing required field: matricula",
				"form missing required field: curso",
				"invalid form value for periodo: type(int) value(abc)",
			)
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
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			HasMessage("error validating form for schema register").
			TraceContains(
				"invalid form value for matricula: type(string) value(2.0221100033e+10)",
				"form missing required field: curso",
				"form missing required field: periodo",
			)
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
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			HasMessage("error validating form for schema register").
			TraceContains("invalid form value for periodo: type(int) value(4.5)")
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
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			HasMessage("error validating form for schema register").
			TraceContains("invalid form value for ativo: type(bool) value(true)")
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
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			HasMessage("error validating form for schema register").
			TraceContains("invalid form value for ativo: type(bool) value(1)")
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
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			HasMessage("error validating form for schema register").
			TraceContains(
				"invalid form value for matricula: type(string) value(true)",
				"form missing required field: curso",
				"form missing required field: periodo",
			)
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
			HasErrID(errx.AuthEmailAlreadyUsed).
			HasMessage("email already in use")
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
				"iss": AsString{projectID, AnyUUID{}},
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
								"fields": map[string]interface{}{
									"curso":     "Ciência da Computação",
									"matricula": "20221100033",
									"periodo":   4,
								},
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

		spec := map[string]interface{}{
			"id": StoreString{Into: &schemaID, Matcher: AnyUUID{}},
		}
		Validate(t, data, spec)

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
				HasErrID(errx.ProjectUserRegisterOnSchemaNoVersion).
				HasMessage("can't register on a schema that has no published version")
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
				HasErrID(errx.ProjectUserRegisterOnSchemaDraft).
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
				HasErrID(errx.ProjectUserRegisterOnSchemaDraft).
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
	})

	// ============================================
	// AMPLIFIED TESTS - Extended Coverage
	// ============================================

	amplifiedClient := suite.NewClient(t)
	amplifiedUser := amplifiedClient.WithCredentials("amplified_schemas@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("Amplified Schema Testing")

	ampProjectID := amplifiedUser.projectID
	var ampSchemaID string
	var ampVersionID string

	t.Run("AmplifiedDraftSchema", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(amplifiedUser.auth)
		data := authClient.POST("/projects/" + ampProjectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "Registration with Options",
				"flow_id":     "registration-v2",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

		spec := map[string]interface{}{
			"id": StoreString{Into: &ampSchemaID, Matcher: AnyUUID{}},
		}
		Validate(t, data, spec)
	})

	t.Run("AmplifiedDraftVersion", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(amplifiedUser.auth)
		data := authClient.POST("/projects/" + ampProjectID + "/schemas/" + ampSchemaID + "/versions/draft").
			Expect(http.StatusCreated).
			RequireDataValue()

		spec := map[string]interface{}{
			"id": StoreString{Into: &ampVersionID, Matcher: AnyUUID{}},
		}
		Validate(t, data, spec)
	})

	t.Run("AmplifiedCreateFieldsWithOptionsAndRules", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(amplifiedUser.auth)
		authClient.POST("/projects/" + ampProjectID + "/schemas/" + ampSchemaID + "/v1").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":         "name",
						"type":        "string",
						"owner":       "user",
						"title":       "Nome Completo",
						"description": "Seu nome completo",
						"required":    true,
						"mutable":     true,
						"position":    0,
					},
					map[string]interface{}{
						"key":      "email",
						"type":     "email",
						"owner":    "user",
						"title":    "Email",
						"required": true,
						"mutable":  true,
						"position": 1,
					},
					// Select field with options
					map[string]interface{}{
						"key":      "user_type",
						"type":     "select",
						"owner":    "user",
						"title":    "Tipo de Usuário",
						"required": true,
						"mutable":  true,
						"position": 2,
						"options": []interface{}{
							map[string]interface{}{"value": "student", "label": "Estudante", "position": 0},
							map[string]interface{}{"value": "teacher", "label": "Professor", "position": 1},
							map[string]interface{}{"value": "staff", "label": "Funcionário", "position": 2},
						},
					},
					// Bool field that other fields depend on
					map[string]interface{}{
						"key":      "is_active",
						"type":     "bool",
						"owner":    "user",
						"title":    "Está Ativo?",
						"required": false,
						"mutable":  true,
						"position": 3,
						"default_value": func() *json.RawMessage {
							b := json.RawMessage("true")
							return &b
						}(),
					},
					// Field with visibility rule (only visible if user_type == "student")
					// NOTE: required=false, so visible but optional
					map[string]interface{}{
						"key":      "student_id",
						"type":     "string",
						"owner":    "user",
						"title":    "RA do Aluno",
						"required": false,
						"mutable":  true,
						"position": 4,
						"visibility_rules": []interface{}{
							map[string]interface{}{
								"depends_on_field_key": "user_type",
								"operator":             "equals",
								"value":                "student",
							},
						},
					},
					// Field with required rule (required if is_active == true)
					map[string]interface{}{
						"key":      "activation_date",
						"type":     "string",
						"owner":    "user",
						"title":    "Data de Ativação",
						"required": false,
						"mutable":  true,
						"position": 5,
						"required_rules": []interface{}{
							map[string]interface{}{
								"depends_on_field_key": "is_active",
								"operator":             "equals",
								"value":                true,
							},
						},
					},
					// Radio field with options
					map[string]interface{}{
						"key":      "shift",
						"type":     "radio",
						"owner":    "user",
						"title":    "Turno",
						"required": false,
						"mutable":  true,
						"position": 6,
						"options": []interface{}{
							map[string]interface{}{"value": "morning", "label": "Matutino", "position": 0},
							map[string]interface{}{"value": "afternoon", "label": "Vespertino", "position": 1},
							map[string]interface{}{"value": "night", "label": "Noturno", "position": 2},
						},
					},
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("created fields")
	})

	t.Run("AmplifiedPublish", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(amplifiedUser.auth)
		authClient.POST("/projects/" + ampProjectID + "/schemas/" + ampSchemaID + "/versions/publish").
			Expect(http.StatusOK)
		authClient.POST("/projects/" + ampProjectID + "/schemas/" + ampSchemaID + "/publish").
			Expect(http.StatusOK)
	})

	// Options validation tests
	t.Run("AmplifiedInvalidOptionValue", func(t *testing.T) {
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "invalid_option@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "John Doe",
					"email":     "john@test.com",
					"user_type": "invalid_type",
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			TraceContains("invalid form value for user_type: type(select) value(invalid_type)")
	})

	t.Run("AmplifiedCaseSensitiveOption", func(t *testing.T) {
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "case_sensitive@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "John Doe",
					"email":     "john@test.com",
					"user_type": "STUDENT",
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			TraceContains("invalid form value for user_type: type(select) value(STUDENT)")
	})

	t.Run("AmplifiedValidOptionSuccess", func(t *testing.T) {
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "valid_option@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "John Doe",
					"email":     "john@test.com",
					"user_type": "student",
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
	})

	// Visibility rules tests
	t.Run("AmplifiedVisibilityRuleHiddenFieldMissing", func(t *testing.T) {
		// user_type = "teacher", so student_id is hidden (not visible)
		// Hidden fields are ignored, so missing student_id is OK
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "hidden_field@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "Jane Doe",
					"email":     "jane@test.com",
					"user_type": "teacher",
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
	})

	t.Run("AmplifiedVisibilityRuleVisibleButOptional", func(t *testing.T) {
		// user_type = "student", student_id is visible
		// But student_id has required=false, so it's optional even when visible
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "visible_optional@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "Student Optional",
					"email":     "optional@test.com",
					"user_type": "student",
					// student_id is visible but optional (required=false)
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
	})

	// Required rules tests
	t.Run("AmplifiedRequiredRuleBasic", func(t *testing.T) {
		// is_active=true triggers required rule for activation_date
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "req_rule_triggered@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "Active No Date",
					"email":     "activenodate@test.com",
					"user_type": "staff",
					"is_active": true,
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			TraceContains("form missing required field: activation_date")
	})

	t.Run("AmplifiedRequiredRuleSatisfied", func(t *testing.T) {
		// is_active=true + activation_date provided = OK
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "req_rule_satisfied@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":            "Active With Date",
					"email":           "activewithdate@test.com",
					"user_type":       "staff",
					"is_active":       true,
					"activation_date": "2024-01-15",
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
	})

	t.Run("AmplifiedRequiredRuleNotTriggered", func(t *testing.T) {
		// is_active=false, required rule NOT triggered
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "req_rule_not_triggered@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "Inactive User",
					"email":     "inactive@test.com",
					"user_type": "staff",
					"is_active": false,
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
	})

	t.Run("AmplifiedRequiredRuleBothMissing", func(t *testing.T) {
		// Both is_active and activation_date missing = OK
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "both_missing@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "Partial User",
					"email":     "partial@test.com",
					"user_type": "teacher",
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
	})

	// Radio field tests
	t.Run("AmplifiedRadioInvalidValue", func(t *testing.T) {
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "radio_invalid@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "Radio User",
					"email":     "radio@test.com",
					"user_type": "student",
					"shift":     "weekend",
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			TraceContains("invalid form value for shift: type(radio) value(weekend)")
	})

	t.Run("AmplifiedRadioValidValue", func(t *testing.T) {
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "radio_valid@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "Radio User",
					"email":     "radio2@test.com",
					"user_type": "teacher",
					"shift":     "night",
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
	})

	// Email validation tests
	t.Run("AmplifiedEmailInvalidFormat", func(t *testing.T) {
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "valid@example.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "Invalid Email User",
					"email":     "not-an-email",
					"user_type": "student",
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			TraceContains("invalid form value for email: type(email) value(not-an-email)")
	})

	t.Run("AmplifiedEmailMissingAt", func(t *testing.T) {
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "valid2@example.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":      "No At User",
					"email":     "invalidemail.com",
					"user_type": "teacher",
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister).
			TraceContains("invalid form value for email: type(email) value(invalidemail.com)")
	})

	t.Run("AmplifiedEmailValidSuccess", func(t *testing.T) {
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "valid_email_test@example.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":       "Valid Email User",
					"email":      "user@company.com",
					"user_type":  "student",
					"student_id": "2023005001",
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
	})

	// Edge cases
	t.Run("AmplifiedExtraFieldsIgnored", func(t *testing.T) {
		// Unknown fields should be silently ignored
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "extra_ignored@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":          "Extra Fields",
					"email":         "extra@test.com",
					"user_type":     "teacher",
					"unknown_field": "should be ignored",
					"another_extra": 123,
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
	})

	t.Run("AmplifiedNullOptionalFields", func(t *testing.T) {
		// Null values for optional fields = treated as "not provided"
		amplifiedClient.POST("/projects/"+ampProjectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "registration-v2").
			WithBody(map[string]interface{}{
				"email":    "null_optional@test.com",
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"name":            "Null User",
					"email":           "null@test.com",
					"user_type":       "staff",
					"is_active":       nil,
					"activation_date": nil,
					"shift":           nil,
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")
	})
}
