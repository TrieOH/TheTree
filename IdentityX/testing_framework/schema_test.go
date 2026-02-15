package testing

import (
	"GoAuth/internal/errx"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func testSchemas(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("schemas@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("schema testing")

	projectID := user.projectID
	rid, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("Couldn't generate uuid for test: %v", err)
	}

	t.Run("PublishSchemaRandomID", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + rid.String() + "/publish").
			Expect(http.StatusUnauthorized).
			HasErrID(errx.SchemaNotOwnedByPrincipal).
			HasMessage("cannot publish a schema you don't own")
	})

	var schemaID string
	t.Run("Draft", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "scti-register-flow",
				"flow_id":     "scti-register",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":                 StoreString{Into: &schemaID, Matcher: AnyUUID{}},
			"project_id":         AsString{projectID, AnyUUID{}},
			"title":              "scti-register-flow",
			"flow_id":            "scti-register",
			"type":               "context",
			"status":             "draft",
			"current_version_id": nil,
			"created_at":         AnyDate{},
			"updated_at":         AnyDate{},
		}

		Validate(t, data, spec)
	})

	t.Run("DraftAnother", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "eenge",
				"flow_id":     "estudante",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":                 AnyUUID{},
			"project_id":         AsString{projectID, AnyUUID{}},
			"title":              "eenge",
			"flow_id":            "estudante",
			"type":               "context",
			"status":             "draft",
			"current_version_id": nil,
			"created_at":         AnyDate{},
			"updated_at":         AnyDate{},
		}

		Validate(t, data, spec)
	})

	t.Run("DraftSameFlowIDAndType", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "eenge",
				"flow_id":     "estudante",
			}).
			Expect(http.StatusConflict).
			HasErrID(errx.SchemaFlowIDAlreadyExistsInType).
			HasMessage("schema with this flow ID already exists in this type")
	})

	t.Run("DraftReservedFlowID", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "Reserved",
				"flow_id":     "none",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.SchemaFlowIDIsReserved).
			HasMessage("flow id can't be the reserved keyword 'none'")
	})

	t.Run("DraftFlowIDSameAsType", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "SameAsType",
				"flow_id":     "context",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.SchemaInvalidFlowID).
			HasMessage("flow id can't be the same as a schema type")
	})

	t.Run("DraftValidation", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)

		t.Run("InvalidType", func(t *testing.T) {
			authClient.POST("/projects/" + projectID + "/schemas").
				WithBody(map[string]interface{}{
					"schema_type": "invalid-type",
					"title":       "test",
					"flow_id":     "test",
				}).
				Expect(http.StatusBadRequest).
				HasErrID(errx.RequestValidationError).
				ValidationError("schema_type must be one of: core, context, sub-context")
		})

		t.Run("FlowIDTooLong", func(t *testing.T) {
			longFlowID := "this-flow-id-is-way-too-long-and-should-fail-validation-because-it-exceeds-63-chars"
			authClient.POST("/projects/" + projectID + "/schemas").
				WithBody(map[string]interface{}{
					"schema_type": "context",
					"title":       "test",
					"flow_id":     longFlowID,
				}).
				Expect(http.StatusBadRequest).
				HasErrID(errx.RequestValidationError).
				ValidationError("flow_id must be at most 63 characters long")
		})
	})

	t.Run("PublishSchemaNoVersion", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusBadRequest).
			HasErrID(errx.SCHEMANoPublishedVersion).
			HasMessage("cannot publish a schema with no versions")
	})

	t.Run("PublishVersionNoDraft", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusUnauthorized).
			HasErrID(errx.SchemaVersionDraftDoesntExist).
			HasMessage("cannot publish a schema with a version draft that doesn't exist")
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
			"title":              "scti-register-flow",
			"flow_id":            "scti-register",
			"type":               "context",
			"status":             "draft",
			"current_version_id": AsString{schemaVersion1ID, AnyUUID{}},
		}

		Validate(t, data, spec)
	})

	t.Run("PublishSchemaDraftVersion", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusBadRequest).
			HasErrID(errx.SchemaHasOnlyDraftVersion).
			HasMessage("cannot publish a schema with only draft versions")
	})

	t.Run("DraftVersionError", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusBadRequest).
			HasErrID(errx.SchemaVersionDraftOnNonPublished).
			HasMessage("new versions can only be drafted from published versions")
	})

	t.Run("PublishVersionFieldsError", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusBadRequest).
			HasErrID(errx.SchemaVersionPublishWithNoFields).
			HasMessage("cannot publish a schema version with no fields")
	})

	t.Run("CreateFieldsSamePosition", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
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
			HasErrID(errx.FIELDSamePositionForMultipleFields).
			HasMessage("two fields can't occupy the same position")
	})

	t.Run("CreateFieldsSameKey", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
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
						"key":         "matricula",
						"type":        "string",
						"owner":       "user",
						"title":       "Numero da Matrícula",
						"description": "Sua matrícula da UENF como aparece no sistema acadêmico",
						"placeholder": "20223200045",
						"required":    true,
						"mutable":     true,
						"position":    1,
					},
				},
			}).
			Expect(http.StatusConflict).
			HasErrID(errx.FIELDSameKeyForMultipleFields).
			HasMessage("two fields can't have the same key")
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
		}

		Validate(t, data, spec)
	})

	t.Run("PublishVersionSuccess", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusOK).
			HasMessage("published schema version")
	})

	t.Run("PublishVersionAlreadyPublished", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusUnauthorized).
			HasErrID(errx.SchemaVersionTryingToPublishPublished).
			HasMessage("cannot publish a schema version that is already published")
	})

	t.Run("PublishSchemaSuccess", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusOK).
			HasMessage("published schema")
	})

	t.Run("PublishSchemaAlreadyPublished", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusUnauthorized).
			HasErrID(errx.SchemaTryingToPublishPublished).
			HasMessage("cannot publish a schema that is already published")
	})

	var schemaVersion2ID string
	t.Run("DraftVersion2", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusCreated).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":             StoreString{Into: &schemaVersion2ID, Matcher: AnyUUID{}},
			"schema_id":      AsString{schemaID, AnyUUID{}},
			"version_number": 2,
		}

		Validate(t, data, spec)
	})

	t.Run("CheckSchemaVersionAfterV2Draft", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.GET("/projects/" + projectID + "/schemas/" + schemaID).
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":                 AsString{schemaID, AnyUUID{}},
			"project_id":         AsString{projectID, AnyUUID{}},
			"title":              "scti-register-flow",
			"flow_id":            "scti-register",
			"type":               "context",
			"status":             "published",
			"current_version_id": AsString{schemaVersion1ID, AnyUUID{}},
		}

		Validate(t, data, spec)
	})

	t.Run("PublishVersion2NoChanges", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusBadRequest).
			HasErrID(errx.SchemaVersionNoChanges).
			HasMessage("cannot publish a version with no changes")
	})

	t.Run("AddFieldToV2FailKeyCheck", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
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
			HasErrID(errx.FIELDInvalidCharactersInKey).
			HasMessage("field key must start with a lowercase letter and contain only lowercase letters, numbers, or underscores")
	})

	t.Run("AddFieldToV2Success", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
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
			HasMessage("created fields").
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"object_id":   AnyUUID{},
				"id":          AnyUUID{},
				"key":         "periodo",
				"type":        "int",
				"owner":       "user",
				"title":       "Período Atual",
				"description": "O período da sua matéria mais avançada da grade",
				"required":    true,
				"mutable":     true,
				"position":    2,
			},
		}

		Validate(t, data, spec)
	})

	t.Run("CreateFieldDuplicateInherited", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/v2").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":         "matricula", // Inherited from v1
						"type":        "string",
						"owner":       "user",
						"title":       "Numero da Matrícula",
						"description": "Duplicate",
						"required":    true,
						"mutable":     true,
						"position":    3, // Different position
					},
				},
			}).
			Expect(http.StatusConflict).
			HasErrID(errx.FIELDSameKeyForMultipleFields).
			HasMessage("two fields can't have the same key")
	})

	t.Run("CreateFieldDuplicateInDraft", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/v2").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":         "periodo", // Created in this draft
						"type":        "int",
						"owner":       "user",
						"title":       "Duplicate",
						"description": "Duplicate",
						"required":    true,
						"mutable":     true,
						"position":    4, // Different position
					},
				},
			}).
			Expect(http.StatusConflict).
			HasErrID(errx.FIELDSameKeyForMultipleFields).
			HasMessage("two fields can't have the same key")
	})

	t.Run("PublishVersion2Success", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusOK).
			HasMessage("published schema version")
	})

	t.Run("GetSchemaVerbose", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		schema := authClient.GET("/projects/" + projectID + "/schemas/" + schemaID + "/verbose").
			Expect(http.StatusOK).
			RequireDataValue()

		// Capture field IDs for cross-version stability checks
		var (
			matriculaV1ID, matriculaV2ID string
			cursoV1ID, cursoV2ID         string
		)

		spec := map[string]interface{}{
			"id":                 AsString{schemaID, AnyUUID{}},
			"project_id":         AsString{projectID, AnyUUID{}},
			"title":              "scti-register-flow",
			"flow_id":            "scti-register",
			"type":               "context",
			"status":             "published",
			"current_version_id": AsString{schemaVersion2ID, AnyUUID{}},
			"created_at":         AnyDate{},
			"updated_at":         AnyDate{},
			"versions": InOrder{
				Specs: []interface{}{
					// Version 2 (newest first in response)
					map[string]interface{}{
						"id":             AsString{schemaVersion2ID, AnyUUID{}},
						"schema_id":      AsString{schemaID, AnyUUID{}},
						"version_number": 2,
						"fields": ByKey{
							Key: "key",
							Spec: map[string]interface{}{
								"matricula": map[string]interface{}{
									"object_id":   AnyUUID{},
									"id":          StoreString{&matriculaV2ID, AnyUUID{}},
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
								"curso": map[string]interface{}{
									"object_id":   AnyUUID{},
									"id":          StoreString{&cursoV2ID, AnyUUID{}},
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
								"periodo": map[string]interface{}{
									"object_id":   AnyUUID{},
									"id":          AnyUUID{},
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
						},
					},
					// Version 1
					map[string]interface{}{
						"id":             AsString{schemaVersion1ID, AnyUUID{}},
						"schema_id":      AsString{schemaID, AnyUUID{}},
						"version_number": 1,
						"fields": ByKey{
							Key: "key",
							Spec: map[string]interface{}{
								"matricula": map[string]interface{}{
									"object_id":   AnyUUID{},
									"id":          StoreString{&matriculaV1ID, AnyUUID{}},
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
								"curso": map[string]interface{}{
									"object_id":   AnyUUID{},
									"id":          StoreString{&cursoV1ID, AnyUUID{}},
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
						},
					},
				},
			},
		}

		Validate(t, schema, spec)

		// Cross-version field ID stability checks
		// Field IDs must match between versions
		require.Equal(t, matriculaV1ID, matriculaV2ID, "matricula field ID changed between versions")
		require.Equal(t, cursoV1ID, cursoV2ID, "curso field ID changed between versions")
	})

	var schemaVersion3ID string
	t.Run("DraftVersion3WithOptionsAndRules", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusCreated).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":                  StoreString{Into: &schemaVersion3ID, Matcher: AnyUUID{}},
			"schema_id":           AsString{schemaID, AnyUUID{}},
			"version_number":      3,
			"status":              "draft",
			"based_on_version_id": AsString{schemaVersion2ID, AnyUUID{}},
			"created_at":          AnyDate{},
			"updated_at":          AnyDate{},
		}

		ValidateExact(t, data, spec)
	})

	t.Run("AddFieldsWithOptionsAndRules", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)

		data := authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/v3").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":      "user_type",
						"type":     "select",
						"owner":    "user",
						"title":    "Tipo de Usuário",
						"required": true,
						"mutable":  true,
						"position": 3,
						"options": []interface{}{
							map[string]interface{}{"value": "student", "label": "Estudante", "position": 0},
							map[string]interface{}{"value": "professor", "label": "Professor", "position": 1},
							map[string]interface{}{"value": "visitor", "label": "Visitante", "position": 2},
						},
					},
					map[string]interface{}{
						"key":      "needs_scholarship",
						"type":     "bool",
						"owner":    "user",
						"title":    "Necessita de Bolsa?",
						"required": true,
						"mutable":  true,
						"position": 4,
					},
					map[string]interface{}{
						"key":         "income",
						"type":        "int",
						"owner":       "user",
						"title":       "Renda Familiar",
						"description": "Renda mensal familiar em reais",
						"required":    false,
						"mutable":     true,
						"position":    5,
						"visibility_rules": []interface{}{
							map[string]interface{}{
								"depends_on_field_key": "needs_scholarship",
								"operator":             "equals",
								"value":                true,
							},
						},
						"required_rules": []interface{}{
							map[string]interface{}{
								"depends_on_field_key": "needs_scholarship",
								"operator":             "equals",
								"value":                true,
							},
						},
					},
					map[string]interface{}{
						"key":      "scholarship_type",
						"type":     "radio",
						"owner":    "user",
						"title":    "Tipo de Bolsa",
						"required": false,
						"mutable":  true,
						"position": 6,
						"options": []interface{}{
							map[string]interface{}{"value": "full", "label": "Integral", "position": 0},
							map[string]interface{}{"value": "partial", "label": "Parcial", "position": 1},
						},
						"visibility_rules": []interface{}{
							map[string]interface{}{
								"depends_on_field_key": "user_type",
								"operator":             "equals",
								"value":                "student",
							},
							map[string]interface{}{
								"depends_on_field_key": "needs_scholarship",
								"operator":             "equals",
								"value":                true,
							},
						},
					},
				},
			}).
			Expect(http.StatusCreated).
			HasMessage("created fields").
			RequireDataValue()

		// The response doesn't include options/rules in the field creation endpoint
		// Only verify core fields exist
		spec := []interface{}{
			map[string]interface{}{
				"object_id":         AnyUUID{},
				"id":                AnyUUID{},
				"schema_id":         AsString{schemaID, AnyUUID{}},
				"schema_version_id": AsString{schemaVersion3ID, AnyUUID{}},
				"key":               "user_type",
				"type":              "select",
				"owner":             "user",
				"title":             "Tipo de Usuário",
				"description":       nil,
				"placeholder":       nil,
				"required":          true,
				"mutable":           true,
				"default_value":     nil,
				"position":          3,
				"created_at":        AnyDate{},
				"updated_at":        AnyDate{},
				// Note: options/rules not returned in create response, only in form/get
			},
			map[string]interface{}{
				"object_id":         AnyUUID{},
				"id":                AnyUUID{},
				"schema_id":         AsString{schemaID, AnyUUID{}},
				"schema_version_id": AsString{schemaVersion3ID, AnyUUID{}},
				"key":               "needs_scholarship",
				"type":              "bool",
				"owner":             "user",
				"title":             "Necessita de Bolsa?",
				"description":       nil,
				"placeholder":       nil,
				"required":          true,
				"mutable":           true,
				"default_value":     nil,
				"position":          4,
				"created_at":        AnyDate{},
				"updated_at":        AnyDate{},
			},
			map[string]interface{}{
				"object_id":         AnyUUID{},
				"id":                AnyUUID{},
				"schema_id":         AsString{schemaID, AnyUUID{}},
				"schema_version_id": AsString{schemaVersion3ID, AnyUUID{}},
				"key":               "income",
				"type":              "int",
				"owner":             "user",
				"title":             "Renda Familiar",
				"description":       "Renda mensal familiar em reais",
				"placeholder":       nil,
				"required":          false,
				"mutable":           true,
				"default_value":     nil,
				"position":          5,
				"created_at":        AnyDate{},
				"updated_at":        AnyDate{},
			},
			map[string]interface{}{
				"object_id":         AnyUUID{},
				"id":                AnyUUID{},
				"schema_id":         AsString{schemaID, AnyUUID{}},
				"schema_version_id": AsString{schemaVersion3ID, AnyUUID{}},
				"key":               "scholarship_type",
				"type":              "radio",
				"owner":             "user",
				"title":             "Tipo de Bolsa",
				"description":       nil,
				"placeholder":       nil,
				"required":          false,
				"mutable":           true,
				"default_value":     nil,
				"position":          6,
				"created_at":        AnyDate{},
				"updated_at":        AnyDate{},
			},
		}

		Validate(t, data, spec)
	})

	t.Run("PublishVersion3Success", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusOK).
			HasMessage("published schema version")
	})

	t.Run("GetLatestFormByID", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.GET("/projects/" + projectID + "/schemas/" + schemaID + "/latest").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":             AsString{schemaVersion3ID, AnyUUID{}},
			"schema_id":      AsString{schemaID, AnyUUID{}},
			"title":          "scti-register-flow",
			"flow_id":        "scti-register",
			"schema_type":    "context",
			"version_id":     AsString{schemaVersion3ID, AnyUUID{}},
			"version_number": 3,
			"status":         "published",
			"created_at":     AnyDate{},
			"updated_at":     AnyDate{},
			"fields": []interface{}{
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "matricula",
					"type":             "string",
					"owner":            "user",
					"title":            "Numero da Matrícula",
					"description":      "Sua matrícula da UENF como aparece no sistema acadêmico",
					"placeholder":      "20223200045",
					"required":         true,
					"mutable":          true,
					"default_value":    nil,
					"position":         0,
					"created_at":       AnyDate{},
					"updated_at":       AnyDate{},
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "curso",
					"type":             "string",
					"owner":            "user",
					"title":            "Curso de Matrícula",
					"description":      "O curso que você está matrículado na UENF",
					"placeholder":      "Ciência da Computação",
					"required":         true,
					"mutable":          true,
					"default_value":    nil,
					"position":         1,
					"created_at":       AnyDate{},
					"updated_at":       AnyDate{},
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "periodo",
					"type":             "int",
					"owner":            "user",
					"title":            "Período Atual",
					"description":      "O período da sua matéria mais avançada da grade",
					"placeholder":      nil,
					"required":         true,
					"mutable":          true,
					"default_value":    nil,
					"position":         2,
					"created_at":       AnyDate{},
					"updated_at":       AnyDate{},
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":            AnyUUID{},
					"object_id":     AnyUUID{},
					"key":           "user_type",
					"type":          "select",
					"owner":         "user",
					"title":         "Tipo de Usuário",
					"description":   nil,
					"placeholder":   nil,
					"required":      true,
					"mutable":       true,
					"default_value": nil,
					"position":      3,
					"created_at":    AnyDate{},
					"updated_at":    AnyDate{},
					"options": []interface{}{
						map[string]interface{}{
							"id":       AnyUUID{},
							"value":    "student",
							"label":    "Estudante",
							"position": 0,
						},
						map[string]interface{}{
							"id":       AnyUUID{},
							"value":    "professor",
							"label":    "Professor",
							"position": 1,
						},
						map[string]interface{}{
							"id":       AnyUUID{},
							"value":    "visitor",
							"label":    "Visitante",
							"position": 2,
						},
					},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "needs_scholarship",
					"type":             "bool",
					"owner":            "user",
					"title":            "Necessita de Bolsa?",
					"description":      nil,
					"placeholder":      nil,
					"required":         true,
					"mutable":          true,
					"default_value":    nil,
					"position":         4,
					"created_at":       AnyDate{},
					"updated_at":       AnyDate{},
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":            AnyUUID{},
					"object_id":     AnyUUID{},
					"key":           "income",
					"type":          "int",
					"owner":         "user",
					"title":         "Renda Familiar",
					"description":   "Renda mensal familiar em reais",
					"placeholder":   nil,
					"required":      false,
					"mutable":       true,
					"default_value": nil,
					"position":      5,
					"created_at":    AnyDate{},
					"updated_at":    AnyDate{},
					"options":       []interface{}{},
					"visibility_rules": []interface{}{
						map[string]interface{}{
							"id":                  AnyUUID{},
							"depends_on_field_id": AnyUUID{},
							"operator":            "equals",
							"value":               true,
						},
					},
					"required_rules": []interface{}{
						map[string]interface{}{
							"id":                  AnyUUID{},
							"depends_on_field_id": AnyUUID{},
							"operator":            "equals",
							"value":               true,
						},
					},
				},
				map[string]interface{}{
					"id":            AnyUUID{},
					"object_id":     AnyUUID{},
					"key":           "scholarship_type",
					"type":          "radio",
					"owner":         "user",
					"title":         "Tipo de Bolsa",
					"description":   nil,
					"placeholder":   nil,
					"required":      false,
					"mutable":       true,
					"default_value": nil,
					"position":      6,
					"created_at":    AnyDate{},
					"updated_at":    AnyDate{},
					"options": []interface{}{
						map[string]interface{}{
							"id":       AnyUUID{},
							"value":    "full",
							"label":    "Integral",
							"position": 0,
						},
						map[string]interface{}{
							"id":       AnyUUID{},
							"value":    "partial",
							"label":    "Parcial",
							"position": 1,
						},
					},
					"visibility_rules": []interface{}{
						map[string]interface{}{
							"id":                  AnyUUID{},
							"depends_on_field_id": AnyUUID{},
							"operator":            "equals",
							"value":               "student",
						},
						map[string]interface{}{
							"id":                  AnyUUID{},
							"depends_on_field_id": AnyUUID{},
							"operator":            "equals",
							"value":               true,
						},
					},
					"required_rules": []interface{}{},
				},
			},
		}

		ValidateExact(t, data, spec)
	})

	t.Run("GetSpecificFormV2", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.GET("/projects/" + projectID + "/schemas/" + schemaID + "/v2").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":             AsString{schemaVersion2ID, AnyUUID{}},
			"schema_id":      AsString{schemaID, AnyUUID{}},
			"title":          "scti-register-flow",
			"flow_id":        "scti-register",
			"schema_type":    "context",
			"version_id":     AsString{schemaVersion2ID, AnyUUID{}},
			"version_number": 2,
			"status":         "published",
			"created_at":     AnyDate{},
			"updated_at":     AnyDate{},
			"fields": []interface{}{
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "matricula",
					"type":             "string",
					"owner":            "user",
					"title":            "Numero da Matrícula",
					"description":      "Sua matrícula da UENF como aparece no sistema acadêmico",
					"placeholder":      "20223200045",
					"required":         true,
					"mutable":          true,
					"position":         0,
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "curso",
					"type":             "string",
					"owner":            "user",
					"title":            "Curso de Matrícula",
					"description":      "O curso que você está matrículado na UENF",
					"placeholder":      "Ciência da Computação",
					"required":         true,
					"mutable":          true,
					"position":         1,
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "periodo",
					"type":             "int",
					"owner":            "user",
					"title":            "Período Atual",
					"description":      "O período da sua matéria mais avançada da grade",
					"placeholder":      nil,
					"required":         true,
					"mutable":          true,
					"position":         2,
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
			},
		}

		Validate(t, data, spec)
	})

	t.Run("GetLatestFormByFlowLookup", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		data := authClient.GET("/projects/"+projectID+"/schemas/lookup/latest").
			WithQuery("flow_id", "scti-register").
			WithQuery("schema_type", "context").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":             AsString{schemaVersion3ID, AnyUUID{}},
			"schema_id":      AsString{schemaID, AnyUUID{}},
			"title":          "scti-register-flow",
			"flow_id":        "scti-register",
			"schema_type":    "context",
			"version_number": 3,
			"status":         "published",
			"created_at":     AnyDate{},
			"updated_at":     AnyDate{},
			"fields": []interface{}{
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "matricula",
					"type":             "string",
					"owner":            "user",
					"title":            "Numero da Matrícula",
					"description":      "Sua matrícula da UENF como aparece no sistema acadêmico",
					"placeholder":      "20223200045",
					"required":         true,
					"mutable":          true,
					"default_value":    nil,
					"position":         0,
					"created_at":       AnyDate{},
					"updated_at":       AnyDate{},
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "curso",
					"type":             "string",
					"owner":            "user",
					"title":            "Curso de Matrícula",
					"description":      "O curso que você está matrículado na UENF",
					"placeholder":      "Ciência da Computação",
					"required":         true,
					"mutable":          true,
					"default_value":    nil,
					"position":         1,
					"created_at":       AnyDate{},
					"updated_at":       AnyDate{},
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "periodo",
					"type":             "int",
					"owner":            "user",
					"title":            "Período Atual",
					"description":      "O período da sua matéria mais avançada da grade",
					"placeholder":      nil,
					"required":         true,
					"mutable":          true,
					"default_value":    nil,
					"position":         2,
					"created_at":       AnyDate{},
					"updated_at":       AnyDate{},
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":            AnyUUID{},
					"object_id":     AnyUUID{},
					"key":           "user_type",
					"type":          "select",
					"owner":         "user",
					"title":         "Tipo de Usuário",
					"description":   nil,
					"placeholder":   nil,
					"required":      true,
					"mutable":       true,
					"default_value": nil,
					"position":      3,
					"created_at":    AnyDate{},
					"updated_at":    AnyDate{},
					"options": []interface{}{
						map[string]interface{}{
							"id":       AnyUUID{},
							"value":    "student",
							"label":    "Estudante",
							"position": 0,
						},
						map[string]interface{}{
							"id":       AnyUUID{},
							"value":    "professor",
							"label":    "Professor",
							"position": 1,
						},
						map[string]interface{}{
							"id":       AnyUUID{},
							"value":    "visitor",
							"label":    "Visitante",
							"position": 2,
						},
					},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":               AnyUUID{},
					"object_id":        AnyUUID{},
					"key":              "needs_scholarship",
					"type":             "bool",
					"owner":            "user",
					"title":            "Necessita de Bolsa?",
					"description":      nil,
					"placeholder":      nil,
					"required":         true,
					"mutable":          true,
					"default_value":    nil,
					"position":         4,
					"created_at":       AnyDate{},
					"updated_at":       AnyDate{},
					"options":          []interface{}{},
					"visibility_rules": []interface{}{},
					"required_rules":   []interface{}{},
				},
				map[string]interface{}{
					"id":            AnyUUID{},
					"object_id":     AnyUUID{},
					"key":           "income",
					"type":          "int",
					"owner":         "user",
					"title":         "Renda Familiar",
					"description":   "Renda mensal familiar em reais",
					"placeholder":   nil,
					"required":      false,
					"mutable":       true,
					"default_value": nil,
					"position":      5,
					"created_at":    AnyDate{},
					"updated_at":    AnyDate{},
					"options":       []interface{}{},
					"visibility_rules": []interface{}{
						map[string]interface{}{
							"id":                  AnyUUID{},
							"depends_on_field_id": AnyUUID{},
							"operator":            "equals",
							"value":               true,
						},
					},
					"required_rules": []interface{}{
						map[string]interface{}{
							"id":                  AnyUUID{},
							"depends_on_field_id": AnyUUID{},
							"operator":            "equals",
							"value":               true,
						},
					},
				},
				map[string]interface{}{
					"id":            AnyUUID{},
					"object_id":     AnyUUID{},
					"key":           "scholarship_type",
					"type":          "radio",
					"owner":         "user",
					"title":         "Tipo de Bolsa",
					"description":   nil,
					"placeholder":   nil,
					"required":      false,
					"mutable":       true,
					"default_value": nil,
					"position":      6,
					"created_at":    AnyDate{},
					"updated_at":    AnyDate{},
					"options": []interface{}{
						map[string]interface{}{
							"id":       AnyUUID{},
							"value":    "full",
							"label":    "Integral",
							"position": 0,
						},
						map[string]interface{}{
							"id":       AnyUUID{},
							"value":    "partial",
							"label":    "Parcial",
							"position": 1,
						},
					},
					"visibility_rules": []interface{}{
						map[string]interface{}{
							"id":                  AnyUUID{},
							"depends_on_field_id": AnyUUID{},
							"operator":            "equals",
							"value":               "student",
						},
						map[string]interface{}{
							"id":                  AnyUUID{},
							"depends_on_field_id": AnyUUID{},
							"operator":            "equals",
							"value":               true,
						},
					},
					"required_rules": []interface{}{},
				},
			},
		}

		Validate(t, data, spec)
	})

	t.Run("GetFormByFlowLookupInvalidType", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.GET("/projects/"+projectID+"/schemas/lookup/latest").
			WithQuery("flow_id", "scti-register").
			WithQuery("schema_type", "core").
			Expect(http.StatusNotFound).
			HasErrID(errx.SQLNotFound).
			HasMessage("schema not found")
	})

	t.Run("GetFormByFlowLookupMissingRequiredQuery", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.GET("/projects/"+projectID+"/schemas/lookup/latest").
			WithQuery("schema_type", "context").
			Expect(http.StatusBadRequest).
			HasErrID(errx.RequestMissingQueryParam).
			HasMessage("missing query parameter: flow_id")
	})

	t.Run("ProjectUserAccessDenied", func(t *testing.T) {
		// Register a project user
		projUser := client.WithCredentials("proj-user-schema@mail.com", ValidPassword).
			ProjectRegister(user.projectID).
			ProjectLogin(user.projectID)

		authClient := suite.NewClient(t).WithAuth(projUser.auth)

		// Try Draft
		authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "forbidden",
				"flow_id":     "forbidden",
			}).
			Expect(http.StatusForbidden).
			HasErrID(errx.AuthNotClient).
			HasMessage("only clients can access this endpoint")

		// Try Publish
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusForbidden).
			HasErrID(errx.AuthNotClient).
			HasMessage("only clients can access this endpoint")

		// Try Get
		authClient.GET("/projects/" + projectID + "/schemas/" + schemaID).
			Expect(http.StatusForbidden).
			HasErrID(errx.AuthNotClient).
			HasMessage("only clients can access this endpoint")

		// Try GetVerbose
		authClient.GET("/projects/" + projectID + "/schemas/" + schemaID + "/verbose").
			Expect(http.StatusForbidden).
			HasErrID(errx.AuthNotClient).
			HasMessage("only clients can access this endpoint")

		// Try Draft Version
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusForbidden).
			HasErrID(errx.AuthNotClient).
			HasMessage("only clients can access this endpoint")
	})
}
