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
			HasErrID(apierr.DBUniqueViolation).
			HasMessage("resource already exists")
	})

	var roleNoDescriptionID string
	t.Run("CreateRoleNoDescription", func(t *testing.T) {
		val := suite.NewClient(t).WithAuth(user.auth).POST("/projects/" + projectID + "/roles").
			WithBody(map[string]interface{}{
				"name": "user",
			}).
			Expect(http.StatusCreated).
			HasMessage("Role Created").
			RequireDataValue()

		spec := map[string]interface{}{
			"id":          StoreString{Into: &roleNoDescriptionID, Matcher: AnyUUID{}},
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
				"id":          StoreString{Into: &roleNoDescriptionID, Matcher: AnyUUID{}},
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
			HasErrID(apierr.RequestValidationError).
			HasMessage("Validation failed").
			TraceContains("name is required")
	})

	t.Run("ForbiddenQueryParam", func(t *testing.T) {
		suite.NewClient(t).WithAuth(user.auth).GET("/projects/"+projectID+"/roles/search").
			WithQuery("something_else", "should_fail").
			Expect(http.StatusBadRequest).
			HasErrID(apierr.RequestUnknownQueryParam).
			HasMessage("unknown query parameter: something_else")
	})
}
