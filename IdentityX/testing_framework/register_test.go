package testing

import (
	"GoAuth/internal/apierr"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func testRegister(t *testing.T, suite *TestSuite) {
	t.Run("Validation", func(t *testing.T) {
		for _, spec := range ValidationTests {
			spec := spec // capture range variable
			t.Run(spec.Name, func(t *testing.T) {
				client := suite.NewClient(t)
				client.POST("/auth/register").
					WithBody(map[string]string{
						"email":    spec.Email,
						"password": spec.Pass,
					}).
					Expect(http.StatusBadRequest).
					HasErrID(apierr.RequestValidationError).
					ValidationError(spec.Errors...)
			})
		}
	})

	t.Run("WeakPasswords", func(t *testing.T) {
		for i, spec := range WeakPasswordTests {
			spec := spec // capture range variable
			i := i       // capture range variable
			t.Run(spec.Name, func(t *testing.T) {
				client := suite.NewClient(t)
				client.POST("/auth/register").
					WithBody(map[string]string{
						"email":    fmt.Sprintf("weak%d@mail.com", i),
						"password": spec.Password,
					}).
					Expect(http.StatusBadRequest).
					HasErrID(apierr.RequestValidationError).
					ValidationError(spec.Errors...)
			})
		}
	})

	t.Run("PasswordTooLong", func(t *testing.T) {
		client := suite.NewClient(t)
		longPass := "A1@" + strings.Repeat("a", 70) // 3 chars + 70 chars = 73 chars > 72
		client.POST("/auth/register").
			WithBody(map[string]string{
				"email":    "longpass@mail.com",
				"password": longPass,
			}).
			Expect(http.StatusBadRequest).
			HasErrID(apierr.RequestValidationError).
			ValidationError("password must be at most 72 characters long")
	})

	// FIXME: Current validation package does not consider trailing spaces as a valid email, fix it and start using the trailing spaces email
	t.Run("EmailNormalization", func(t *testing.T) {
		client := suite.NewClient(t)
		// email := "  MixedCase@Example.Com  "
		email := "MixedCase@Example.Com"
		normalizedEmail := "mixedcase@example.com"
		password := ValidPassword

		// Register with messy email
		client.POST("/auth/register").
			WithBody(map[string]string{
				"email":    email,
				"password": password,
			}).
			Expect(http.StatusCreated).
			HasMessage("Registered user")

		// Login with normalized email
		client.POST("/auth/login").
			WithBody(map[string]string{
				"email":    normalizedEmail,
				"password": password,
			}).
			Expect(http.StatusOK).
			HasMessage("Logged in")

		// Login with messy email again (normalization happens on input)
		client.POST("/auth/login").
			WithBody(map[string]string{
				"email":    email,
				"password": password,
			}).
			Expect(http.StatusOK).
			HasMessage("Logged in")

		// Try to register again with normalized email - should conflict
		client.POST("/auth/register").
			WithBody(map[string]string{
				"email":    normalizedEmail,
				"password": password,
			}).
			Expect(http.StatusConflict).
			HasErrID(apierr.AuthEmailAlreadyUsed).
			HasMessage("email already in use")
	})

	t.Run("Success", func(t *testing.T) {
		client := suite.NewClient(t)
		client.WithCredentials("new@mail.com", ValidPassword).Register()
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		client := suite.NewClient(t)
		email := "duplicate@mail.com"
		client.WithCredentials(email, ValidPassword).Register()

		client.POST("/auth/register").
			WithBody(map[string]string{
				"email":    email,
				"password": ValidPassword,
			}).
			Expect(http.StatusConflict).
			HasErrID(apierr.AuthEmailAlreadyUsed).
			HasMessage("email already in use")
	})
}
