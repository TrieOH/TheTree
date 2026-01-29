package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"
)

func testPermissions(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("permissions@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("permissions testing")

	projectID := user.projectID
	var permissionID string
	t.Run("CreatePermission", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "event:*",
				"action":     "create",
				"conditions": nil,
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         StoreString{Into: &permissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "event:*",
			"action":     "create",
			"conditions": nil,
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("GetPermissionByID", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/" + projectID + "/permissions/" + permissionID).
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         StoreString{Into: &permissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "event:*",
			"action":     "create",
			"conditions": nil,
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	var anotherPermissionID string
	t.Run("CreateAnotherPermission", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "event:123/activity:*",
				"action":     "attendance:mark",
				"conditions": nil,
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         StoreString{Into: &anotherPermissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "event:123/activity:*",
			"action":     "attendance:mark",
			"conditions": nil,
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("ListProjectPermissions", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/" + projectID + "/permissions").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: anotherPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:123/activity:*",
				"action":     "attendance:mark",
				"conditions": nil,
				"created_at": AnyDate{},
			},
			map[string]interface{}{
				"id":         AsString{Value: permissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:*",
				"action":     "create",
				"conditions": nil,
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("GetPermissionByObject", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/"+projectID+"/permissions").
			WithQuery("object", "event:123/activity:*").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: anotherPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:123/activity:*",
				"action":     "attendance:mark",
				"conditions": nil,
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("GetPermissionByAction", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/"+projectID+"/permissions").
			WithQuery("action", "create").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: permissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:*",
				"action":     "create",
				"conditions": nil,
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	var createProductPermissionID string
	t.Run("CreateCreateProductPermission", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "event:123/product:*",
				"action":     "create",
				"conditions": nil,
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         StoreString{Into: &createProductPermissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "event:123/product:*",
			"action":     "create",
			"conditions": nil,
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("GetPermissionByActionAgain", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/"+projectID+"/permissions").
			WithQuery("action", "create").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: createProductPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:123/product:*",
				"action":     "create",
				"conditions": nil,
				"created_at": AnyDate{},
			},
			map[string]interface{}{
				"id":         AsString{Value: permissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:*",
				"action":     "create",
				"conditions": nil,
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	var editProductPermissionID string
	t.Run("CreateEditProductPermission", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "event:123/product:*",
				"action":     "edit",
				"conditions": nil,
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         StoreString{Into: &editProductPermissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "event:123/product:*",
			"action":     "edit",
			"conditions": nil,
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("GetPermissionByObjectAgain", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/"+projectID+"/permissions").
			WithQuery("object", "event:123/product:*").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: editProductPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:123/product:*",
				"action":     "edit",
				"conditions": nil,
				"created_at": AnyDate{},
			},
			map[string]interface{}{
				"id":         AsString{Value: createProductPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:123/product:*",
				"action":     "create",
				"conditions": nil,
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("GetPermissionByObjectAndAction", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/"+projectID+"/permissions").
			WithQuery("object", "event:123/product:*").
			WithQuery("action", "edit").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: editProductPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:123/product:*",
				"action":     "edit",
				"conditions": nil,
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("CreateEventMasterAdminPermission", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "event:456",
				"action":     "*",
				"conditions": nil,
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         AnyUUID{},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "event:456",
			"action":     "*",
			"conditions": nil,
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("CreateMasterUserPermission", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "*",
				"action":     "*",
				"conditions": nil,
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         AnyUUID{},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "*",
			"action":     "*",
			"conditions": nil,
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("CreateDuplicatePermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "event:*",
				"action":     "create",
				"conditions": nil,
			}).
			Expect(http.StatusConflict).
			HasErrID(apierr.DBUniqueViolation).
			HasMessage("resource already exists")
	})

	t.Run("CreatePermissionNoAction", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "event:*",
				"action":     "",
				"conditions": nil,
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.RequestValidationError).
			HasMessage("Validation failed").
			TraceContains("action is required")
	})

	t.Run("CreatePermissionNoObject", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "",
				"action":     "create",
				"conditions": nil,
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.RequestValidationError).
			HasMessage("Validation failed").
			TraceContains("object is required")
	})

	t.Run("CreateEmptyPermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/"+projectID+"/permissions").
			WithBody(map[string]interface{}{
				"object":     "",
				"action":     "",
				"conditions": nil,
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.RequestValidationError).
			HasMessage("Validation failed").
			TraceContains("object is required", "action is required")
	})

	t.Run("CreateInvalidObjectPermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "bogus-value",
				"action":     "attendance:mark",
				"conditions": nil,
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.PermissionInvalidObject).
			HasMessage("invalid permission object: bogus-value")
	})

	t.Run("CreateInvalidActionPermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "event:*",
				"action":     "what:**",
				"conditions": nil,
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.PermissionInvalidAction).
			HasMessage("invalid permission action: what:*")
	})

	t.Run("NotAllowedQueryParam", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/permissions").
			WithQuery("not-allowed", "should_deny").
			Expect(http.StatusBadRequest).
			HasErrID(apierr.RequestUnknownQueryParam).
			HasMessage("unknown query parameter: not-allowed")
	})
}
