package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"
)

func testLogin(t *testing.T, suite *TestSuite) {
	// Create user in parent test context
	client := suite.NewClient(t)
	user := client.WithCredentials("login@mail.com", ValidPassword).Register()

	t.Run("WrongPassword", func(t *testing.T) {
		client := suite.NewClient(t)
		client.POST("/auth/login").
			WithBody(map[string]string{
				"email":    user.email,
				"password": "WrongPass123!",
			}).
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.AuthInvalidCredentials).
			HasMessage("invalid email or password")
	})

	t.Run("WrongEmail", func(t *testing.T) {
		client := suite.NewClient(t)
		client.POST("/auth/login").
			WithBody(map[string]string{
				"email":    "wrong@mail.com",
				"password": user.password,
			}).
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.AuthInvalidCredentials).
			HasMessage("invalid email or password")
	})

	t.Run("Success", func(t *testing.T) {
		client := suite.NewClient(t)
		client.WithCredentials(user.email, user.password).Login()
	})

	t.Run("Logout", func(t *testing.T) {
		client := suite.NewClient(t)
		loggedInUser := client.WithCredentials(user.email, user.password).Login()
		loggedInUser.Logout()

		// Try using revoked session
		loggedInUser.POST("/auth/logout").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SessionRevoked).
			HasMessage("session not found or revoked")
	})
}
