package testing

import (
	"fmt"
	"net/http"
	"testing"
)

func testRegister(t *testing.T, suite *TestSuite) {
	t.Run("Validation", func(t *testing.T) {
		for _, spec := range ValidationTests {
			spec := spec // capture range variable
			t.Run(spec.Name, func(t *testing.T) {
				client := suite.Client(t)
				client.POST("/auth/register").
					WithBody(map[string]string{
						"email":    spec.Email,
						"password": spec.Pass,
					}).
					Expect(http.StatusBadRequest).
					ValidationError(spec.Errors...)
			})
		}
	})

	t.Run("WeakPasswords", func(t *testing.T) {
		for i, spec := range WeakPasswordTests {
			spec := spec // capture range variable
			i := i       // capture range variable
			t.Run(spec.Name, func(t *testing.T) {
				client := suite.Client(t)
				client.POST("/auth/register").
					WithBody(map[string]string{
						"email":    fmt.Sprintf("weak%d@mail.com", i),
						"password": spec.Password,
					}).
					Expect(http.StatusBadRequest).
					ValidationError(spec.Errors...)
			})
		}
	})

	t.Run("Success", func(t *testing.T) {
		client := suite.Client(t)
		client.User("new@mail.com", ValidPassword).Register()
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		client := suite.Client(t)
		email := "duplicate@mail.com"
		client.User(email, ValidPassword).Register()

		client.POST("/auth/register").
			WithBody(map[string]string{
				"email":    email,
				"password": ValidPassword,
			}).
			Expect(http.StatusConflict).
			Error("go-auth-test", "error registering user")
	})
}
