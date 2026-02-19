package testing

import (
	"GoAuth/internal/errx"
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
				"object": "document",
				"action": "create",
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         StoreString{Into: &permissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "document",
			"action":     "create",
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
			"object":     "document",
			"action":     "create",
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	var anotherPermissionID string
	t.Run("CreateAnotherPermission", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "event",
				"action": "read",
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         StoreString{Into: &anotherPermissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "event",
			"action":     "read",
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
				"object":     "event",
				"action":     "read",
				"created_at": AnyDate{},
			},
			map[string]interface{}{
				"id":         AsString{Value: permissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "document",
				"action":     "create",
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("GetPermissionByObject", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/"+projectID+"/permissions").
			WithQuery("object", "event").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: anotherPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event",
				"action":     "read",
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
				"object":     "document",
				"action":     "create",
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
				"object": "product",
				"action": "create",
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         StoreString{Into: &createProductPermissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "product",
			"action":     "create",
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
				"object":     "product",
				"action":     "create",
				"created_at": AnyDate{},
			},
			map[string]interface{}{
				"id":         AsString{Value: permissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "document",
				"action":     "create",
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
				"object": "product",
				"action": "edit",
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         StoreString{Into: &editProductPermissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "product",
			"action":     "edit",
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("GetPermissionByObjectAgain", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/"+projectID+"/permissions").
			WithQuery("object", "product").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: editProductPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "product",
				"action":     "edit",
				"created_at": AnyDate{},
			},
			map[string]interface{}{
				"id":         AsString{Value: createProductPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "product",
				"action":     "create",
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("GetPermissionByObjectAndAction", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/"+projectID+"/permissions").
			WithQuery("object", "product").
			WithQuery("action", "edit").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: editProductPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "product",
				"action":     "edit",
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("CreateAdminPermission", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "user",
				"action": "delete",
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         AnyUUID{},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "user",
			"action":     "delete",
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("CreateMasterPermission", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "system",
				"action": "admin",
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         AnyUUID{},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "system",
			"action":     "admin",
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("CreateDuplicatePermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "document",
				"action": "create",
			}).
			Expect(http.StatusConflict).
			HasErrID(errx.PERMissionAlreadyExists).
			HasMessage("permission with object(document) and action(create) already exists")
	})

	t.Run("CreatePermissionNoAction", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "document",
				"action": "",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.RequestValidationError).
			HasMessage("Validation failed").
			TraceContains("action is required")
	})

	t.Run("CreatePermissionNoObject", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "",
				"action": "create",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.RequestValidationError).
			HasMessage("Validation failed").
			TraceContains("object is required")
	})

	t.Run("CreateEmptyPermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/"+projectID+"/permissions").
			WithBody(map[string]interface{}{
				"object": "",
				"action": "",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.RequestValidationError).
			HasMessage("Validation failed").
			TraceContains("object is required", "action is required")
	})

	t.Run("CreateInvalidObjectPermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "event:*",
				"action": "read",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.PERMissionInvalidObject).
			HasMessage("invalid permission object: (event:*)")
	})

	t.Run("CreateInvalidActionPermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "document",
				"action": "attendance:mark",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.PERMissionInvalidAction).
			HasMessage("invalid permission action: (attendance:mark)")
	})

	t.Run("CreateWildcardObjectPermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "*",
				"action": "read",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.PERMissionInvalidObject).
			HasMessage("invalid permission object: (*)")
	})

	t.Run("CreateWildcardActionPermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object": "document",
				"action": "*",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(errx.PERMissionInvalidAction).
			HasMessage("invalid permission action: (*)")
	})

	t.Run("NotAllowedQueryParam", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/permissions").
			WithQuery("not-allowed", "should_deny").
			Expect(http.StatusBadRequest).
			HasErrID(errx.RequestUnknownQueryParam).
			HasMessage("unknown query parameter: not-allowed")
	})
}
