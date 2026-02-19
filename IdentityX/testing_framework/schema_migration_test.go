package testing

import (
	"GoAuth/internal/errx"
	"context"
	"fmt"
	"net/http"
	"testing"
)

func testSchemaMigration(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	// Create a client user (project owner)
	owner := client.WithCredentials("migration_owner@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("migration project")

	projectID := owner.ProjectID()
	authClient := suite.NewClient(t).WithAuth(owner.auth)

	var schemaID string
	var version1ID string

	t.Run("SetupSchemaV1", func(t *testing.T) {
		// 1. Create Schema
		data := authClient.WithT(t).POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "Migration Schema",
				"flow_id":     "migrate",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

		spec := map[string]interface{}{
			"id": StoreString{Into: &schemaID, Matcher: AnyUUID{}},
		}
		Validate(t, data, spec)

		// 2. Create Version 1 Draft
		data = authClient.WithT(t).POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusCreated).
			RequireDataValue()

		spec = map[string]interface{}{
			"id": StoreString{Into: &version1ID, Matcher: AnyUUID{}},
		}
		Validate(t, data, spec)

		// 3. Add Field 1 (Required)
		authClient.WithT(t).POST("/projects/" + projectID + "/schemas/" + schemaID + "/v1").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":      "field1",
						"type":     "string",
						"owner":    "user",
						"title":    "Field 1",
						"position": 0,
						"required": true,
					},
				},
			}).
			Expect(http.StatusCreated)

		// 4. Publish Version 1
		authClient.WithT(t).POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusOK)

		// 5. Publish Schema
		authClient.WithT(t).POST("/projects/" + projectID + "/schemas/" + schemaID + "/publish").
			Expect(http.StatusOK)
	})

	// Create a project user and register with V1
	userEmail := "migrated_user@mail.com"
	t.Run("RegisterUserV1", func(t *testing.T) {
		client.WithT(t).WithCredentials(userEmail, ValidPassword).
			POST("/projects/"+projectID+"/register").
			WithQuery("schema_type", "context").
			WithQuery("flow_id", "migrate").
			WithBody(map[string]interface{}{
				"email":    userEmail,
				"password": ValidPassword,
				"custom_fields": map[string]interface{}{
					"field1": "value1",
				},
			}).
			Expect(http.StatusCreated)
	})

	var userClient *Client
	var entityID string
	t.Run("CheckUserIsUpToDateV1", func(t *testing.T) {
		loginResp := client.WithT(t).WithCredentials(userEmail, ValidPassword).
			POST("/projects/" + projectID + "/login").
			WithBody(map[string]interface{}{
				"email":    userEmail,
				"password": ValidPassword,
			}).
			Expect(http.StatusOK)

		data := loginResp.RequireDataValue()
		Validate(t, data, map[string]interface{}{
			"is_up_to_date": true,
		})

		auth := loginResp.AuthCookies()
		userClient = suite.NewClient(t).WithCredentials(userEmail, ValidPassword).WithAuth(auth)

		// Get user entity ID using framework
		meData := userClient.WithT(t).GET("/sessions/me").Expect(http.StatusOK).RequireDataValue()
		Validate(t, meData, map[string]interface{}{
			"access": map[string]interface{}{
				"sub": map[string]interface{}{
					"id": StoreString{Into: &entityID, Matcher: AnyUUID{}},
				},
			},
		})
	})

	t.Run("VerifyRedisEntryCreatedAfterCompatibilityCheck", func(t *testing.T) {
		cacheKey := fmt.Sprintf("compat:%s:%s:%s", projectID, version1ID, entityID)
		val, err := suite.Redis.Get(context.Background(), cacheKey).Result()
		if err != nil {
			t.Fatalf("Expected key %s to exist in Redis, but got error: %v", cacheKey, err)
		}
		if val != "true" {
			t.Errorf("Expected Redis value for key %s to be 'true', but got '%s'", cacheKey, val)
		}
	})

	var permissionID string
	t.Run("SetupPermissions", func(t *testing.T) {
		// Give user a permission
		data := authClient.WithT(t).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "document",
				"action": "read",
			}).
			Expect(http.StatusCreated).
			RequireDataValue()

		spec := map[string]interface{}{
			"id": StoreString{Into: &permissionID, Matcher: AnyUUID{}},
		}
		Validate(t, data, spec)

		authClient.WithT(t).POST("/projects/" + projectID + "/identities/" + entityID + "/permissions").
			WithBody(map[string]interface{}{
				"permission_id": permissionID,
			}).
			Expect(http.StatusOK)
	})

	t.Run("PermissionWorksV1", func(t *testing.T) {
		data := authClient.WithT(t).POST("/authz/check").
			WithBody(map[string]interface{}{
				"entity_id":  entityID,
				"object":     "document",
				"action":     "read",
				"project_id": projectID,
			}).
			Expect(http.StatusOK).
			RequireDataValue()

		Validate(t, data, map[string]interface{}{
			"allowed": true,
		})
	})

	var version2ID string
	t.Run("UpgradeSchemaToV2", func(t *testing.T) {
		// 1. Create Version 2 Draft
		data := authClient.WithT(t).POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/draft").
			Expect(http.StatusCreated).
			RequireDataValue()

		spec := map[string]interface{}{
			"id": StoreString{Into: &version2ID, Matcher: AnyUUID{}},
		}
		Validate(t, data, spec)

		// 2. Add Field 2 (Required) - This makes V1 users outdated
		authClient.WithT(t).POST("/projects/" + projectID + "/schemas/" + schemaID + "/v2").
			WithBody(map[string]interface{}{
				"fields": []interface{}{
					map[string]interface{}{
						"key":      "field2",
						"type":     "string",
						"owner":    "user",
						"title":    "Field 2",
						"position": 1,
						"required": true,
					},
				},
			}).
			Expect(http.StatusCreated)

		// 3. Publish Version 2
		authClient.WithT(t).POST("/projects/" + projectID + "/schemas/" + schemaID + "/versions/publish").
			Expect(http.StatusOK)

		// 4. Verify Redis cleared for this project and old version
		cacheKeyV1 := fmt.Sprintf("compat:%s:%s:%s", projectID, version1ID, entityID)
		_, err := suite.Redis.Get(context.Background(), cacheKeyV1).Result()
		if err == nil {
			t.Errorf("Expected key %s to be cleared from Redis after publishing new version, but it still exists", cacheKeyV1)
		}
	})

	t.Run("CheckUserIsOutdatedV2", func(t *testing.T) {
		// Login should now return is_up_to_date: false
		data := client.WithT(t).WithCredentials(userEmail, ValidPassword).
			POST("/projects/" + projectID + "/login").
			WithBody(map[string]interface{}{
				"email":    userEmail,
				"password": ValidPassword,
			}).
			Expect(http.StatusOK).
			RequireDataValue()

		Validate(t, data, map[string]interface{}{
			"is_up_to_date": false,
		})

		// 2. Verify Redis entry for V2 is created (and it is 'false')
		cacheKeyV2 := fmt.Sprintf("compat:%s:%s:%s", projectID, version2ID, entityID)
		val, err := suite.Redis.Get(context.Background(), cacheKeyV2).Result()
		if err != nil {
			t.Fatalf("Expected key %s to exist in Redis, but got error: %v", cacheKeyV2, err)
		}
		if val != "false" {
			t.Errorf("Expected Redis value for key %s to be 'false', but got '%s'", cacheKeyV2, val)
		}
	})

	t.Run("PermissionBlockedV2", func(t *testing.T) {
		// Permission check should now fail with UserSchemaOutdated
		authClient.WithT(t).POST("/authz/check").
			WithBody(map[string]interface{}{
				"entity_id":  entityID,
				"object":     "document",
				"action":     "read",
				"project_id": projectID,
			}).
			Expect(http.StatusForbidden).
			HasErrID(errx.AuthUserSchemaOutdated)
	})

	t.Run("GetUpgradeForm", func(t *testing.T) {
		data := userClient.WithT(t).GET("/projects/" + projectID + "/upgrade-form").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"SchemaID":      AsString{Value: schemaID, Matcher: AnyUUID{}},
				"VersionID":     AsString{Value: version2ID, Matcher: AnyUUID{}},
				"VersionNumber": 2,
				"Title":         "Migration Schema",
				"FlowID":        "migrate",
				"SchemaType":    "context",
				"Fields": ByKey{
					Key: "Key",
					Spec: map[string]interface{}{
						"field1": map[string]interface{}{
							"Key":      "field1",
							"Type":     "string",
							"Required": true,
							"Position": 0,
						},
						"field2": map[string]interface{}{
							"Key":      "field2",
							"Type":     "string",
							"Required": true,
							"Position": 1,
						},
					},
				},
			},
		}
		Validate(t, data, spec)
	})

	t.Run("UpdateMetadataFailValidation", func(t *testing.T) {
		// Try to update but miss the required field2
		userClient.WithT(t).POST("/projects/" + projectID + "/metadata").
			WithBody(map[string]interface{}{
				"custom_fields": map[string]interface{}{
					"field1": "value1-updated",
				},
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.FIELDValidationErrorOnSchemaRegister)
	})

	t.Run("UpdateMetadataSuccess", func(t *testing.T) {
		userClient.WithT(t).POST("/projects/" + projectID + "/metadata").
			WithBody(map[string]interface{}{
				"custom_fields": map[string]interface{}{
					"field2": "value2",
				},
			}).
			Expect(http.StatusOK)

		// Verify Redis entry for V2 is now 'true'
		cacheKeyV2 := fmt.Sprintf("compat:%s:%s:%s", projectID, version2ID, entityID)
		val, err := suite.Redis.Get(context.Background(), cacheKeyV2).Result()
		if err != nil {
			t.Fatalf("Expected key %s to exist in Redis, but got error: %v", cacheKeyV2, err)
		}
		if val != "true" {
			t.Errorf("Expected Redis value for key %s to be 'true', but got '%s'", cacheKeyV2, val)
		}
	})

	t.Run("CheckUserIsUpToDateAgain", func(t *testing.T) {
		userClient = userClient.WithT(t).ProjectLogin(projectID)
		data := userClient.POST("/projects/" + projectID + "/login").
			WithBody(map[string]interface{}{
				"email":    userEmail,
				"password": ValidPassword,
			}).
			Expect(http.StatusOK).
			RequireDataValue()

		Validate(t, data, map[string]interface{}{
			"is_up_to_date": true,
		})
	})

	t.Run("PermissionUnblocked", func(t *testing.T) {
		data := authClient.WithT(t).POST("/authz/check").
			WithBody(map[string]interface{}{
				"entity_id":  entityID,
				"object":     "document",
				"action":     "read",
				"project_id": projectID,
			}).
			Expect(http.StatusOK).
			RequireDataValue()

		Validate(t, data, map[string]interface{}{
			"allowed": true,
		})
	})

	t.Run("MetadataReflectedInSession", func(t *testing.T) {
		meData := userClient.WithT(t).GET("/sessions/me").Expect(http.StatusOK).RequireDataValue()

		spec := map[string]interface{}{
			"access": map[string]interface{}{
				"sub": map[string]interface{}{
					"metadata": map[string]interface{}{
						"context": map[string]interface{}{
							"migrate": map[string]interface{}{
								"schema_version_id": version2ID,
								"fields": map[string]interface{}{
									"field1": "value1",
									"field2": "value2",
								},
							},
						},
					},
				},
			},
		}
		Validate(t, meData, spec)
	})
}
