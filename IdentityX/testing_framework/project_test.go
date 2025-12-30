package testing

import (
	"net/http"
	"testing"
)

func testProjects(t *testing.T, suite *TestSuite) {
	client := suite.Client(t)
	user := client.User("projects@mail.com", ValidPassword).Register().Login()

	var projectID string

	t.Run("CreateProject", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects").
			WithBody(map[string]interface{}{
				"project_name": "Test Project",
				"metadata":     map[string]string{"env": "test"},
			}).
			Expect(http.StatusCreated).
			Data()

		projectID = data.Value("id").String().Raw()
		data.Value("project_name").String().IsEqual("Test Project")
	})

	t.Run("ListProjects", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		arr := authClient.GET("/projects").
			Expect(http.StatusOK).
			DataArray()

		arr.Length().IsEqual(1)
		arr.Value(0).Object().Value("id").IsEqual(projectID)
	})

	t.Run("GetProject", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.GET("/projects/" + projectID).
			Expect(http.StatusOK).
			Data()

		data.Value("id").String().IsEqual(projectID)
		data.Value("project_name").String().IsEqual("Test Project")
	})

	t.Run("UpdateProject", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.PATCH("/projects/" + projectID).
			WithBody(map[string]interface{}{
				"project_name": "Updated Project",
				"metadata":     map[string]string{"env": "prod"},
			}).
			Expect(http.StatusOK).
			Data()

		data.Value("project_name").String().IsEqual("Updated Project")
	})

	t.Run("GetProjectJWKS", func(t *testing.T) {
		jwksClient := suite.Client(t)
		obj := jwksClient.GET("/projects/" + projectID + "/.well-known/jwks.json").
			Expect(http.StatusOK).
			JSON()

		obj.Value("keys").Array().NotEmpty()
	})

	t.Run("DeleteProject", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		authClient.DELETE("/projects/"+projectID).
			Expect(http.StatusOK).
			Success("go-auth-test", "Deleted project")

		// Verify deletion
		authClient.GET("/projects").
			Expect(http.StatusOK).
			DataArray().
			Length().IsEqual(0)
	})
}
