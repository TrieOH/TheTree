package testing

import (
	"net/http"
	"testing"
)

func testSchemas(t *testing.T, suite *TestSuite) {
	client := suite.Client(t)
	user := client.User("schemas@mail.com", ValidPassword).
		Register().
		Login().
		CreateProject("schema testing")

	var projectID string = user.ProjectID
	t.Run("DraftSchema", func(t *testing.T) {
		authClient := suite.Client(t).Auth(user.auth)
		data := authClient.POST("/projects/" + projectID + "/schemas").
			WithBody(map[string]interface{}{
				"schema_type": "context",
				"title":       "scti-register-flow",
				"flow_id":     "scti-register",
			}).
			Expect(http.StatusCreated).
			Data()

		data.Value("project_id").String().IsEqual(projectID)
		data.Value("title").String().IsEqual("scti-register-flow")
		data.Value("flow_id").String().IsEqual("scti-register")
		data.Value("id").String().NotEmpty()
		data.Value("type").String().IsEqual("context")
		data.Value("status").String().IsEqual("draft")
		data.Value("current_version_id").IsNull()
	})
}
