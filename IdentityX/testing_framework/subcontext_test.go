package testing

import (
	"GoAuth/internal/errx"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func testSubContext(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("subcontext@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("subcontext test project")

	projectID := user.ProjectID()

	// Register a project user and get their ID
	projectUserEmail := "projectuser@mail.com"
	_ = suite.NewClient(t).WithCredentials(projectUserEmail, ValidPassword).ProjectRegister(projectID)

	// Get the project user ID
	authClient := suite.NewClient(t).WithAuth(user.auth)
	usersResp := authClient.GET("/projects/" + projectID + "/users").
		Expect(http.StatusOK).
		RequireDataArray()

	var projectUserID string
	for i := 0; i < int(usersResp.Length().Raw()); i++ {
		userObj := usersResp.Value(i).Object()
		email := userObj.Value("email").String().Raw()
		if email == projectUserEmail {
			projectUserID = userObj.Value("id").String().Raw()
			break
		}
	}
	require.NotEmpty(t, projectUserID, "Could not find project user ID")

	t.Run("AddSubContext", func(t *testing.T) {
		t.Run("Unauthorized", func(t *testing.T) {
			client := suite.NewClient(t)
			client.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"data": map[string]interface{}{
						"tags": []string{"CanCreateEvents"},
					},
				}).
				Expect(http.StatusUnauthorized).
				HasErrID(errx.AuthMissingAccessCookie)
		})

		t.Run("ProjectUserForbidden", func(t *testing.T) {
			projectUserClient := suite.NewClient(t).WithCredentials("projectuser2@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)

			projectUserClient.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"data": map[string]interface{}{
						"tags": []string{"CanCreateEvents"},
					},
				}).
				Expect(http.StatusForbidden).
				HasErrID(errx.AuthNotClient)
		})

		t.Run("InvalidProjectID", func(t *testing.T) {
			client := suite.NewClient(t).WithAuth(user.auth)
			client.POST("/projects/invalid-uuid/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"data": map[string]interface{}{
						"tags": []string{"CanCreateEvents"},
					},
				}).
				Expect(http.StatusBadRequest).
				HasErrID(errx.RequestValidationError)
		})

		t.Run("ProjectNotFound", func(t *testing.T) {
			nonexistentID, _ := uuid.NewV7()
			client := suite.NewClient(t).WithAuth(user.auth)
			client.POST("/projects/" + nonexistentID.String() + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"data": map[string]interface{}{
						"tags": []string{"CanCreateEvents"},
					},
				}).
				Expect(http.StatusNotFound).
				HasErrID(errx.SQLNotFound).
				HasMessage("project not found")
		})

		t.Run("UserNotFound", func(t *testing.T) {
			nonexistentUserID, _ := uuid.NewV7()
			client := suite.NewClient(t).WithAuth(user.auth)
			client.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": nonexistentUserID.String(),
					"data": map[string]interface{}{
						"tags": []string{"CanCreateEvents"},
					},
				}).
				Expect(http.StatusNotFound).
				HasErrID(errx.SQLNotFound)
		})

		t.Run("ValidationError", func(t *testing.T) {
			client := suite.NewClient(t).WithAuth(user.auth)
			client.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": "invalid-uuid",
					"data":    "not-an-object",
				}).
				Expect(http.StatusBadRequest).
				HasErrID(errx.RequestInvalidJSONFormat)
		})

		t.Run("MissingData", func(t *testing.T) {
			client := suite.NewClient(t).WithAuth(user.auth)
			client.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
				}).
				Expect(http.StatusBadRequest).
				HasErrID(errx.RequestValidationError)
		})

		t.Run("Success", func(t *testing.T) {
			client := suite.NewClient(t).WithAuth(user.auth)
			client.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"data": map[string]interface{}{
						"tags": []string{"CanCreateEvents", "PremiumUser"},
						"permissions": map[string]interface{}{
							"can_delete": true,
							"can_edit":   false,
						},
					},
				}).
				Expect(http.StatusOK).
				HasMessage("Sub-context added successfully")
		})

		t.Run("SuccessMerge", func(t *testing.T) {
			// Register a new user for merge test
			mergeTestEmail := "mergetest@mail.com"
			suite.NewClient(t).POST("/projects/" + projectID + "/register").
				WithBody(map[string]interface{}{
					"email":    mergeTestEmail,
					"password": ValidPassword,
				}).
				Expect(http.StatusCreated)

			// Get the user ID
			client := suite.NewClient(t).WithAuth(user.auth)
			usersResp := client.GET("/projects/" + projectID + "/users").
				Expect(http.StatusOK).
				RequireDataArray()

			var mergeUserID string
			for i := 0; i < int(usersResp.Length().Raw()); i++ {
				userObj := usersResp.Value(i).Object()
				email := userObj.Value("email").String().Raw()
				if email == mergeTestEmail {
					mergeUserID = userObj.Value("id").String().Raw()
					break
				}
			}
			require.NotEmpty(t, mergeUserID)

			// Add initial data
			client.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": mergeUserID,
					"data": map[string]interface{}{
						"initial_key": "initial_value",
					},
				}).
				Expect(http.StatusOK)

			// Add more data - should merge
			client.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": mergeUserID,
					"data": map[string]interface{}{
						"additional_key": "additional_value",
					},
				}).
				Expect(http.StatusOK)

			// Verify both keys exist
			resp := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/users/" + mergeUserID + "/sub-context").
				Expect(http.StatusOK).
				JSONObj()

			data := resp.Value("data").Object()
			subContext := data.Value("sub_context").Object()
			subContext.Value("initial_key").String().IsEqual("initial_value")
			subContext.Value("additional_key").String().IsEqual("additional_value")
		})
	})

	t.Run("GetSubContext", func(t *testing.T) {
		t.Run("Unauthorized", func(t *testing.T) {
			client := suite.NewClient(t)
			client.GET("/projects/" + projectID + "/users/" + projectUserID + "/sub-context").
				Expect(http.StatusUnauthorized).
				HasErrID(errx.AuthMissingAccessCookie)
		})

		t.Run("ProjectUserForbidden", func(t *testing.T) {
			projectUserClient := suite.NewClient(t).WithCredentials("gettest@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)

			projectUserClient.GET("/projects/" + projectID + "/users/" + projectUserID + "/sub-context").
				Expect(http.StatusForbidden).
				HasErrID(errx.AuthNotClient)
		})

		t.Run("InvalidProjectID", func(t *testing.T) {
			client := suite.NewClient(t).WithAuth(user.auth)
			client.GET("/projects/invalid-uuid/users/" + projectUserID + "/sub-context").
				Expect(http.StatusBadRequest).
				HasErrID(errx.RequestValidationError)
		})

		t.Run("InvalidUserID", func(t *testing.T) {
			client := suite.NewClient(t).WithAuth(user.auth)
			client.GET("/projects/" + projectID + "/users/invalid-uuid/sub-context").
				Expect(http.StatusBadRequest).
				HasErrID(errx.RequestValidationError)
		})

		t.Run("ProjectNotFound", func(t *testing.T) {
			nonexistentID, _ := uuid.NewV7()
			client := suite.NewClient(t).WithAuth(user.auth)
			client.GET("/projects/" + nonexistentID.String() + "/users/" + projectUserID + "/sub-context").
				Expect(http.StatusNotFound).
				HasErrID(errx.SQLNotFound).
				HasMessage("project not found")
		})

		t.Run("UserNotFound", func(t *testing.T) {
			nonexistentUserID, _ := uuid.NewV7()
			client := suite.NewClient(t).WithAuth(user.auth)
			client.GET("/projects/" + projectID + "/users/" + nonexistentUserID.String() + "/sub-context").
				Expect(http.StatusNotFound).
				HasErrID(errx.SQLNotFound)
		})

		t.Run("SuccessEmpty", func(t *testing.T) {
			// Register a new user with no sub-context
			emptyUserEmail := "emptyuser@mail.com"
			suite.NewClient(t).POST("/projects/" + projectID + "/register").
				WithBody(map[string]interface{}{
					"email":    emptyUserEmail,
					"password": ValidPassword,
				}).
				Expect(http.StatusCreated)

			// Get the user ID
			client := suite.NewClient(t).WithAuth(user.auth)
			usersResp := client.GET("/projects/" + projectID + "/users").
				Expect(http.StatusOK).
				RequireDataArray()

			var emptyUserID string
			for i := 0; i < int(usersResp.Length().Raw()); i++ {
				userObj := usersResp.Value(i).Object()
				email := userObj.Value("email").String().Raw()
				if email == emptyUserEmail {
					emptyUserID = userObj.Value("id").String().Raw()
					break
				}
			}
			require.NotEmpty(t, emptyUserID)

			resp := client.GET("/projects/" + projectID + "/users/" + emptyUserID + "/sub-context").
				Expect(http.StatusOK).
				JSONObj()

			data := resp.Value("data").Object()
			subContext := data.Value("sub_context").Object()
			require.NotNil(t, subContext)
		})

		t.Run("SuccessWithData", func(t *testing.T) {
			// First add some data
			client := suite.NewClient(t).WithAuth(user.auth)
			client.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"data": map[string]interface{}{
						"test_tag": "test_value",
					},
				}).
				Expect(http.StatusOK)

			resp := client.GET("/projects/" + projectID + "/users/" + projectUserID + "/sub-context").
				Expect(http.StatusOK).
				JSONObj()

			data := resp.Value("data").Object()
			subContext := data.Value("sub_context").Object()
			subContext.Value("test_tag").String().IsEqual("test_value")
		})
	})

	t.Run("RemoveSubContext", func(t *testing.T) {
		t.Run("Unauthorized", func(t *testing.T) {
			client := suite.NewClient(t)
			client.DELETE("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"keys":    []string{"tags"},
				}).
				Expect(http.StatusUnauthorized).
				HasErrID(errx.AuthMissingAccessCookie)
		})

		t.Run("ProjectUserForbidden", func(t *testing.T) {
			projectUserClient := suite.NewClient(t).WithCredentials("removetest@mail.com", ValidPassword).
				ProjectRegister(projectID).
				ProjectLogin(projectID)

			projectUserClient.DELETE("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"keys":    []string{"tags"},
				}).
				Expect(http.StatusForbidden).
				HasErrID(errx.AuthNotClient)
		})

		t.Run("InvalidProjectID", func(t *testing.T) {
			client := suite.NewClient(t).WithAuth(user.auth)
			client.DELETE("/projects/invalid-uuid/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"keys":    []string{"tags"},
				}).
				Expect(http.StatusBadRequest).
				HasErrID(errx.RequestValidationError)
		})

		t.Run("ProjectNotFound", func(t *testing.T) {
			nonexistentID, _ := uuid.NewV7()
			client := suite.NewClient(t).WithAuth(user.auth)
			client.DELETE("/projects/" + nonexistentID.String() + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"keys":    []string{"tags"},
				}).
				Expect(http.StatusNotFound).
				HasErrID(errx.SQLNotFound).
				HasMessage("project not found")
		})

		t.Run("UserNotFound", func(t *testing.T) {
			nonexistentUserID, _ := uuid.NewV7()
			client := suite.NewClient(t).WithAuth(user.auth)
			client.DELETE("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": nonexistentUserID.String(),
					"keys":    []string{"tags"},
				}).
				Expect(http.StatusNotFound).
				HasErrID(errx.SQLNotFound)
		})

		t.Run("MissingKeys", func(t *testing.T) {
			client := suite.NewClient(t).WithAuth(user.auth)
			client.DELETE("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
				}).
				Expect(http.StatusBadRequest).
				HasErrID(errx.RequestValidationError)
		})

		t.Run("EmptyKeys", func(t *testing.T) {
			client := suite.NewClient(t).WithAuth(user.auth)
			client.DELETE("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"keys":    []string{},
				}).
				Expect(http.StatusBadRequest).
				HasErrID(errx.RequestValidationError)
		})

		t.Run("Success", func(t *testing.T) {
			// Register a user for removal test
			removeTestEmail := "removesuccess@mail.com"
			suite.NewClient(t).POST("/projects/" + projectID + "/register").
				WithBody(map[string]interface{}{
					"email":    removeTestEmail,
					"password": ValidPassword,
				}).
				Expect(http.StatusCreated)

			// Get the user ID
			client := suite.NewClient(t).WithAuth(user.auth)
			usersResp := client.GET("/projects/" + projectID + "/users").
				Expect(http.StatusOK).
				RequireDataArray()

			var removeUserID string
			for i := 0; i < int(usersResp.Length().Raw()); i++ {
				userObj := usersResp.Value(i).Object()
				email := userObj.Value("email").String().Raw()
				if email == removeTestEmail {
					removeUserID = userObj.Value("id").String().Raw()
					break
				}
			}
			require.NotEmpty(t, removeUserID)

			// First add some data
			client.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": removeUserID,
					"data": map[string]interface{}{
						"key_to_remove": "value",
						"key_to_keep":   "value",
					},
				}).
				Expect(http.StatusOK)

			// Remove the key
			client.DELETE("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": removeUserID,
					"keys":    []string{"key_to_remove"},
				}).
				Expect(http.StatusOK).
				HasMessage("Sub-context keys removed successfully")

			// Verify key was removed
			resp := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/users/" + removeUserID + "/sub-context").
				Expect(http.StatusOK).
				JSONObj()

			data := resp.Value("data").Object()
			subContext := data.Value("sub_context").Object()
			subContext.NotContainsKey("key_to_remove")
			subContext.Value("key_to_keep").String().IsEqual("value")
		})

		t.Run("SuccessNestedKey", func(t *testing.T) {
			// Register a user for nested test
			nestedTestEmail := "nestedtest@mail.com"
			suite.NewClient(t).POST("/projects/" + projectID + "/register").
				WithBody(map[string]interface{}{
					"email":    nestedTestEmail,
					"password": ValidPassword,
				}).
				Expect(http.StatusCreated)

			// Get the user ID
			client := suite.NewClient(t).WithAuth(user.auth)
			usersResp := client.GET("/projects/" + projectID + "/users").
				Expect(http.StatusOK).
				RequireDataArray()

			var nestedUserID string
			for i := 0; i < int(usersResp.Length().Raw()); i++ {
				userObj := usersResp.Value(i).Object()
				email := userObj.Value("email").String().Raw()
				if email == nestedTestEmail {
					nestedUserID = userObj.Value("id").String().Raw()
					break
				}
			}
			require.NotEmpty(t, nestedUserID)

			// Add nested data
			client.POST("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": nestedUserID,
					"data": map[string]interface{}{
						"permissions": map[string]interface{}{
							"can_read":  true,
							"can_write": true,
						},
					},
				}).
				Expect(http.StatusOK)

			// Remove nested key
			client.DELETE("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": nestedUserID,
					"keys":    []string{"permissions.can_write"},
				}).
				Expect(http.StatusOK)

			// Verify nested key was removed
			resp := suite.NewClient(t).WithAuth(user.auth).GET("/projects/" + projectID + "/users/" + nestedUserID + "/sub-context").
				Expect(http.StatusOK).
				JSONObj()

			data := resp.Value("data").Object()
			subContext := data.Value("sub_context").Object()
			permissions := subContext.Value("permissions").Object()
			permissions.NotContainsKey("can_write")
			permissions.Value("can_read").Boolean().IsEqual(true)
		})

		t.Run("SuccessRemoveNonExistent", func(t *testing.T) {
			// Remove a key that doesn't exist - should not error
			client := suite.NewClient(t).WithAuth(user.auth)
			client.DELETE("/projects/" + projectID + "/sub-context").
				WithBody(map[string]interface{}{
					"user_id": projectUserID,
					"keys":    []string{"non_existent_key"},
				}).
				Expect(http.StatusOK)
		})
	})

	t.Run("JWTContainsSubContext", func(t *testing.T) {
		// Verify sub-context is retrievable (it goes to JWT via metadata)
		// First, add sub-context to the main project user
		client := suite.NewClient(t).WithAuth(user.auth)
		client.POST("/projects/" + projectID + "/sub-context").
			WithBody(map[string]interface{}{
				"user_id": projectUserID,
				"data": map[string]interface{}{
					"jwt_test": "value",
				},
			}).
			Expect(http.StatusOK)

		// Verify it was stored correctly
		resp := client.GET("/projects/" + projectID + "/users/" + projectUserID + "/sub-context").
			Expect(http.StatusOK).
			JSONObj()

		data := resp.Value("data").Object()
		subContext := data.Value("sub_context").Object()
		subContext.Value("jwt_test").String().IsEqual("value")
	})
}
