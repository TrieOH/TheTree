package testing

import (
	"GoAuth/internal/apierr"
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
	rid, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("Couldn't generate random uuid for test: %v", err)
	}

	t.Run("PublishSchemaRandomID", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + rid.String() + "/publish").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SchemaNotOwnedByPrincipal).
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
			HasErrID(apierr.SchemaFlowIDAlreadyExistsInType).
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
			HasErrID(apierr.SchemaFlowIDIsReserved).
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
			HasErrID(apierr.SchemaInvalidFlowID).
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
				HasErrID(apierr.RequestValidationError).
				ValidationError("(schema_type)")
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
				HasErrID(apierr.RequestValidationError).
				ValidationError("(flow_id)")
		})
	})

	t.Run("PublishSchemaNoVersion", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SchemaNoPublishedVersion).
			HasMessage("cannot publish a schema with no versions")
	})

	t.Run("PublishVersionNoDraft", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SchemaVersionDraftDoesntExist).
			HasMessage("cannot publish a schema version draft that doesn't exist")
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

	t.Run("PublishSchemaOnlyDraft", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SchemaHasOnlyDraftVersion).
			HasMessage("cannot publish a schema with only draft versions")
	})

	t.Run("DraftVersionError", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SchemaVersionDraftOnNonPublished).
			HasMessage("new versions can only be drafted from published versions")
	})

	t.Run("PublishVersionFieldsError", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SchemaVersionPublishWithNoFields).
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
			HasErrID(apierr.FieldSamePositionForMultipleFields).
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
			HasErrID(apierr.FieldSameKeyForMultipleFields).
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
			HasErrID(apierr.SchemaVersionTryingToPublishPublished).
			HasMessage("cannot publish a schema version that isn't a draft")
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
			HasErrID(apierr.SchemaTryingToPublishPublished).
			HasMessage("cannot publish a schema that isn't a draft")
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
			"current_version_id": AsString{schemaVersion2ID, AnyUUID{}},
		}

		Validate(t, data, spec)
	})

	t.Run("PublishVersion2NoChanges", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusBadRequest).
			HasErrID(apierr.SchemaVersionNoChanges).
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
			HasErrID(apierr.FieldInvalidCharactersInKey).
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
			HasErrID(apierr.FieldSameKeyForMultipleFields).
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
			HasErrID(apierr.FieldSameKeyForMultipleFields).
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
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.AuthNotClient).
			HasMessage("only clients can access this endpoint")

		// Try Publish
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.AuthNotClient).
			HasMessage("only clients can access this endpoint")

		// Try Get
		authClient.GET("/projects/" + projectID + "/schemas/" + schemaID).
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.AuthNotClient).
			HasMessage("only clients can access this endpoint")

		// Try GetVerbose
		authClient.GET("/projects/" + projectID + "/schemas/" + schemaID + "/verbose").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.AuthNotClient).
			HasMessage("only clients can access this endpoint")

		// Try Draft Version
		authClient.POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.AuthNotClient).
			HasMessage("only clients can access this endpoint")
	})
}
