package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRoles(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("roles@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("roles testing")

	projectID := user.projectID
	var adminRoleID string
	var adminCreateDate string
	t.Run("CreateRole", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.POST("/projects/" + projectID + "/roles").
			WithBody(map[string]interface{}{
				"name":        "admin",
				"description": "can do stuff",
			}).
			Expect(http.StatusCreated).
			HasMessage("Role Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          StoreString{Into: &adminRoleID, Matcher: AnyUUID{}},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "admin",
			"description": "can do stuff",
			"created_at":  AnyDate{},
			"updated_at":  StoreString{Into: &adminCreateDate, Matcher: AnyDate{}},
		}

		Validate(t, val, spec)
	})

	t.Run("UpdateRoleComment", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).PATCH("/projects/" + projectID + "/roles/" + adminRoleID).
			WithBody(map[string]interface{}{
				"description": "can do stuff and more stuff",
			}).
			Expect(http.StatusOK)
	})

	var adminUpdate1Date string
	t.Run("GetRoleByID", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/" + projectID + "/roles/" + adminRoleID).
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          AsString{Value: adminRoleID, Matcher: AnyUUID{}},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "admin",
			"description": "can do stuff and more stuff",
			"created_at":  AnyDate{},
			"updated_at":  StoreString{Into: &adminUpdate1Date, Matcher: AnyDate{}},
		}

		Validate(t, val, spec)

		assert.NotEqual(t, adminCreateDate, adminUpdate1Date)
	})

	t.Run("UpdateRoleCommentAgain", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).PATCH("/projects/" + projectID + "/roles/" + adminRoleID).
			WithBody(map[string]interface{}{
				"description": "can do stuff and more stuff but not that",
			}).
			Expect(http.StatusOK)
	})

	var adminUpdate2Date string
	t.Run("GetRoleByName", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.GET("/projects/"+projectID+"/roles/search").
			WithQuery("name", "admin").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          AsString{Value: adminRoleID, Matcher: AnyUUID{}},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "admin",
			"description": "can do stuff and more stuff but not that",
			"created_at":  AnyDate{},
			"updated_at":  StoreString{Into: &adminUpdate2Date, Matcher: AnyDate{}},
		}

		Validate(t, val, spec)

		assert.NotEqual(t, adminUpdate1Date, adminUpdate2Date)
	})

	t.Run("CreateDuplicateRole", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/roles").
			WithBody(map[string]interface{}{
				"name": "admin",
			}).
			Expect(http.StatusConflict).
			HasErrID(apierr.ID(apierr.ROLENameAlreadyTaken.String())). // FIXME make the error "role name 'NAME' already taken"
			HasMessage("role name already taken")
	})

	var userRoleID string
	t.Run("CreateUserRole", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/roles").
			WithBody(map[string]interface{}{
				"name": "user",
			}).
			Expect(http.StatusCreated).
			HasMessage("Role Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          StoreString{Into: &userRoleID, Matcher: AnyUUID{}},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "user",
			"description": nil,
			"created_at":  AnyDate{},
			"updated_at":  AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("ListProjectRoles", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/roles").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":          StoreString{Into: &userRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "user",
				"description": nil,
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},
			},
			map[string]interface{}{
				"id":          AsString{Value: adminRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "admin",
				"description": "can do stuff and more stuff but not that",
				"created_at":  AnyDate{},
				"updated_at":  AsString{Value: adminUpdate2Date, Matcher: AnyDate{}},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("CreateRoleNoName", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/roles").
			WithBody(map[string]interface{}{
				"name": "",
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.ID(apierr.RequestValidationError.String())).
			HasMessage("Validation failed").
			TraceContains("name is required")
	})

	t.Run("ForbiddenQueryParam", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/roles/search").
			WithQuery("something_else", "should_fail").
			Expect(http.StatusBadRequest).
			HasErrID(apierr.ID(apierr.RequestUnknownQueryParam.String())).
			HasMessage("unknown query parameter: something_else")
	})

	var createEventPermissionID string
	t.Run("CreateEventPermission", func(t *testing.T) {
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
			"id":         StoreString{Into: &createEventPermissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "event:*",
			"action":     "create",
			"conditions": nil,
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	var markAttendancePermissionID string
	t.Run("CreateAttendanceMarkPermission", func(t *testing.T) {
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
			"id":         StoreString{Into: &markAttendancePermissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "event:123/activity:*",
			"action":     "attendance:mark",
			"conditions": nil,
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	var attendActivity321PermissionID string
	t.Run("CreateActivityAttendancePermission", func(t *testing.T) {
		authClient := suite.NewClient(t).WithAuth(user.auth)
		val := authClient.POST("/projects/" + projectID + "/permissions").
			WithBody(map[string]interface{}{
				"object":     "event:123/activity:321",
				"action":     "attend",
				"conditions": nil,
			}).
			Expect(http.StatusCreated).
			HasMessage("Permission Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":         StoreString{Into: &attendActivity321PermissionID, Matcher: AnyUUID{}},
			"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
			"object":     "event:123/activity:321",
			"action":     "attend",
			"conditions": nil,
			"created_at": AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("AddAdminPermissions", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/roles/" + adminRoleID + "/permissions/" + createEventPermissionID).
			Expect(http.StatusOK).
			HasMessage("Added permission to role")

		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/roles/" + adminRoleID + "/permissions/" + markAttendancePermissionID).
			Expect(http.StatusOK).
			HasMessage("Added permission to role")
	})

	t.Run("AddUserPermissions", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/roles/" + userRoleID + "/permissions/" + attendActivity321PermissionID).
			Expect(http.StatusOK).
			HasMessage("Added permission to role")
	})

	t.Run("GetAdminPermissions", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/roles/" + adminRoleID + "/permissions").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: markAttendancePermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:123/activity:*",
				"action":     "attendance:mark",
				"conditions": nil,
				"created_at": AnyDate{},
			},
			map[string]interface{}{
				"id":         AsString{Value: createEventPermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:*",
				"action":     "create",
				"conditions": nil,
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("GetUserPermissions", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/roles/" + userRoleID + "/permissions").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: attendActivity321PermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:123/activity:321",
				"action":     "attend",
				"conditions": nil,
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	t.Run("RemoveAdminPermission", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).DELETE("/projects/" + projectID + "/roles/" + adminRoleID + "/permissions/" + createEventPermissionID).
			Expect(http.StatusOK).
			HasMessage("Removed permission from role")
	})

	t.Run("GetAdminPermissionsAgain", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/roles/" + adminRoleID + "/permissions").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":         AsString{Value: markAttendancePermissionID, Matcher: AnyUUID{}},
				"project_id": AsString{Value: projectID, Matcher: AnyUUID{}},
				"object":     "event:123/activity:*",
				"action":     "attendance:mark",
				"conditions": nil,
				"created_at": AnyDate{},
			},
		}

		Validate(t, val, spec)
	})

	//FIXME check if scope allow multiple id on the same name
	var eventScopeID string
	t.Run("CreateEventScope", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/scopes").
			WithBody(map[string]interface{}{
				"name":        "events",
				"external_id": nil,
			}).
			Expect(http.StatusCreated).
			HasMessage("Scope Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          StoreString{Into: &eventScopeID, Matcher: AnyUUID{}},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "events",
			"external_id": nil,
			"type":        "project_scope",
			"created_at":  AnyDate{},
		}

		Validate(t, val, spec)
	})

	var projectUserID string
	projectUser := client.WithCredentials("roles-user@mail.com", ValidPassword).
		ProjectRegister(projectID).
		ProjectLogin(projectID)

	// Fetch identity_id from sessions list as /sessions/me returns user_id in sub
	meValue := projectUser.GET("/sessions/me").
		Expect(http.StatusOK).
		RequireDataObject()

	projectUserID = meValue.Value("access").Object().Value("sub").Object().Value("id").String().Raw()

	t.Run("GiveUserAdminRole", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  adminRoleID,
				"scope_id": eventScopeID,
			}).
			Expect(http.StatusOK).
			HasMessage("Added role to user")
	})

	t.Run("GetUserRoles", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":          AsString{Value: adminRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "admin",
				"description": "can do stuff and more stuff but not that",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    AsString{Value: eventScopeID, Matcher: AnyUUID{}},
				"scope_name":  "events",
				"external_id": nil,
			},
		}

		Validate(t, val, spec)
	})

	t.Run("GiveUserUserRole", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  userRoleID,
				"scope_id": eventScopeID,
			}).
			Expect(http.StatusOK).
			HasMessage("Added role to user")
	})

	t.Run("GetUserRolesAgain", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":          AsString{Value: adminRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "admin",
				"description": "can do stuff and more stuff but not that",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    AsString{Value: eventScopeID, Matcher: AnyUUID{}},
				"scope_name":  "events",
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: userRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "user",
				"description": nil,
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    AsString{Value: eventScopeID, Matcher: AnyUUID{}},
				"scope_name":  "events",
				"external_id": nil,
			},
		}

		Validate(t, val, spec)
	})

	var scopelessRoleID string
	t.Run("CreateScopelessRole", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/roles").
			WithBody(map[string]interface{}{
				"name":        "scopeless",
				"description": "this role should be project wide",
			}).
			Expect(http.StatusCreated).
			HasMessage("Role Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          StoreString{Into: &scopelessRoleID, Matcher: AnyUUID{}},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "scopeless",
			"description": "this role should be project wide",
			"created_at":  AnyDate{},
			"updated_at":  AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("GiveUserScopelessRole", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  scopelessRoleID,
				"scope_id": nil,
			}).
			Expect(http.StatusOK).
			HasMessage("Added role to user")
	})

	t.Run("GetUserRolesAfterScopeless", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":          AsString{Value: adminRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "admin",
				"description": "can do stuff and more stuff but not that",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    AsString{Value: eventScopeID, Matcher: AnyUUID{}},
				"scope_name":  "events",
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    nil,
				"scope_name":  nil,
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: userRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "user",
				"description": nil,
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    AsString{Value: eventScopeID, Matcher: AnyUUID{}},
				"scope_name":  "events",
				"external_id": nil,
			},
		}

		Validate(t, val, spec)
	})

	t.Run("TakeUserAdminRole", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).DELETE("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  adminRoleID,
				"scope_id": eventScopeID,
			}).
			Expect(http.StatusOK).
			HasMessage("Removed role from user")
	})

	t.Run("GetUserRolesAfterTake", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    nil,
				"scope_name":  nil,
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: userRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "user",
				"description": nil,
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    AsString{Value: eventScopeID, Matcher: AnyUUID{}},
				"scope_name":  "events",
				"external_id": nil,
			},
		}

		Validate(t, val, spec)
	})

	t.Run("GiveUserScopelessRoleOnAScope", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  scopelessRoleID,
				"scope_id": eventScopeID,
			}).
			Expect(http.StatusOK).
			HasMessage("Added role to user")
	})

	t.Run("GetUserRolesAfterScopelessOnScope", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    nil,
				"scope_name":  nil,
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    eventScopeID,
				"scope_name":  "events",
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: userRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "user",
				"description": nil,
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    AsString{Value: eventScopeID, Matcher: AnyUUID{}},
				"scope_name":  "events",
				"external_id": nil,
			},
		}

		Validate(t, val, spec)
	})

	var activityScopeID string
	t.Run("CreateActivityScope", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/scopes").
			WithBody(map[string]interface{}{
				"name": "activities",
			}).
			Expect(http.StatusCreated).
			HasMessage("Scope Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          StoreString{Into: &activityScopeID, Matcher: AnyUUID{}},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "activities",
			"external_id": nil,
			"type":        "project_scope",
			"created_at":  AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("GiveUserScopelessRoleOnActivityScope", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  scopelessRoleID,
				"scope_id": activityScopeID,
			}).
			Expect(http.StatusOK).
			HasMessage("Added role to user")
	})

	t.Run("GetUserRolesAfterScopelessOnActivityScope", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    nil,
				"scope_name":  nil,
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    activityScopeID,
				"scope_name":  "activities",
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    eventScopeID,
				"scope_name":  "events",
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: userRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "user",
				"description": nil,
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    AsString{Value: eventScopeID, Matcher: AnyUUID{}},
				"scope_name":  "events",
				"external_id": nil,
			},
		}

		Validate(t, val, spec)
	})

	var activitySubScopeID string
	t.Run("CreateActivitySubScope", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/scopes").
			WithBody(map[string]interface{}{
				"name":        "activities",
				"external_id": "123",
			}).
			Expect(http.StatusCreated).
			HasMessage("Scope Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          StoreString{Into: &activitySubScopeID, Matcher: AnyUUID{}},
			"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
			"name":        "activities",
			"external_id": "123",
			"type":        "project_scope",
			"created_at":  AnyDate{},
		}

		Validate(t, val, spec)
	})

	t.Run("GiveUserScopelessRoleOnActivitySubScope", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  scopelessRoleID,
				"scope_id": activitySubScopeID,
			}).
			Expect(http.StatusOK).
			HasMessage("Added role to user")
	})

	t.Run("GetUserRolesAfterScopelessOnActivitySubScope", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := []interface{}{
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    nil,
				"scope_name":  nil,
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    activityScopeID,
				"scope_name":  "activities",
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    activitySubScopeID,
				"scope_name":  "activities",
				"external_id": "123",
			},
			map[string]interface{}{
				"id":          AsString{Value: scopelessRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "scopeless",
				"description": "this role should be project wide",
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    eventScopeID,
				"scope_name":  "events",
				"external_id": nil,
			},
			map[string]interface{}{
				"id":          AsString{Value: userRoleID, Matcher: AnyUUID{}},
				"project_id":  AsString{Value: projectID, Matcher: AnyUUID{}},
				"name":        "user",
				"description": nil,
				"created_at":  AnyDate{},
				"updated_at":  AnyDate{},

				"scope_id":    AsString{Value: eventScopeID, Matcher: AnyUUID{}},
				"scope_name":  "events",
				"external_id": nil,
			},
		}

		Validate(t, val, spec)
	})

	t.Run("GiveUserDuplicateScopelessRole", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  scopelessRoleID,
				"scope_id": nil,
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.ID(apierr.ROLEAlreadyGranted.String())).
			HasMessage("scopeless role already granted to user")
	})

	t.Run("GiveUserDuplicateScopelessRoleOnEventScope", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  scopelessRoleID,
				"scope_id": eventScopeID,
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.ID(apierr.ROLEAlreadyGranted.String())).
			HasMessage("scopeless role already granted to user")
	})

	t.Run("GiveUserDuplicateScopelessRoleOnActivityScope", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  scopelessRoleID,
				"scope_id": activityScopeID,
			}).
			Expect(http.StatusConflict).
			HasErrID(apierr.ID(apierr.ROLEAlreadyGranted.String())).
			HasMessage("user already has this role in the specified scope")
	})

	t.Run("GiveUserDuplicateScopelessRoleOnActivitySubScope", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  scopelessRoleID,
				"scope_id": activitySubScopeID,
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.ID(apierr.ROLEAlreadyGranted.String())).
			HasMessage("scopeless role already granted to user")
	})

	// FIXME right now taking a role from a user succeeds with OK because of how DELETE works in SQL
	// FIXME Switch to :execrows and if no row was modified return and error warning user already doesnt have the role
	t.Run("TakeRoleAlreadyTaken", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).DELETE("/projects/" + projectID + "/identities/" + projectUserID + "/roles").
			WithBody(map[string]interface{}{
				"role_id":  adminRoleID,
				"scope_id": eventScopeID,
			}).
			Expect(http.StatusOK).
			HasMessage("Removed role from user")
	})

	t.Run("GiveUserDirectEventPermissionNoScope", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/permissions").
			WithBody(map[string]interface{}{
				"permission_id": createEventPermissionID,
				"scope_id":      nil,
			}).
			Expect(http.StatusOK).
			HasMessage("Added permission to user")
	})

	t.Run("GiveUserDirectEventPermissionWithEventScope", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/permissions").
			WithBody(map[string]interface{}{
				"permission_id": createEventPermissionID,
				"scope_id":      eventScopeID,
			}).
			Expect(http.StatusOK).
			HasMessage("Added permission to user")
	})

	t.Run("GiveUserDirectDuplicateEventPermissionNoScope", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/identities/" + projectUserID + "/permissions").
			WithBody(map[string]interface{}{
				"permission_id": createEventPermissionID,
				"scope_id":      nil,
			}).
			Expect(http.StatusConflict).
			HasErrID(apierr.ID(apierr.PERMissionAlreadyGranted.String())).
			HasMessage("user already has this permission in the specified scope")
	})

	t.Run("TakeDirectEventPermissionWithEventScopeFromUser", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).DELETE("/projects/" + projectID + "/identities/" + projectUserID + "/permissions").
			WithBody(map[string]interface{}{
				"permission_id": createEventPermissionID,
				"scope_id":      eventScopeID,
			}).
			Expect(http.StatusOK).
			HasMessage("Removed permission from user")
	})

	// FIXME right now taking a permission from a user succeeds with OK because of how DELETE works in SQL
	// FIXME Switch to :execrows and if no row was modified return and error warning user already doesnt have the permission
	t.Run("TakeDirectEventPermissionWithEventScopeFromUserAlreadyTaken", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).DELETE("/projects/" + projectID + "/identities/" + projectUserID + "/permissions").
			WithBody(map[string]interface{}{
				"permission_id": createEventPermissionID,
				"scope_id":      eventScopeID,
			}).
			Expect(http.StatusOK).
			HasMessage("Removed permission from user")
	})
}
