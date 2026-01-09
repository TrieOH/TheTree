package testing

import (
	"GoAuth/internal/apierr"
	"net/http"
	"testing"
)

func testLogin(t *testing.T, suite *TestSuite) {
	// Create user in parent test context
	client := suite.NewClient(t)
	user := client.NewUser("login@mail.com", ValidPassword).Register()

	t.Run("WrongPassword", func(t *testing.T) {
		client := suite.NewClient(t)
		client.POST("/auth/login").
			WithBody(map[string]string{
				"email":    user.Email,
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
				"password": user.Password,
			}).
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.AuthInvalidCredentials).
			HasMessage("invalid email or password")
	})

	t.Run("Success", func(t *testing.T) {
		client := suite.NewClient(t)
		client.NewUser(user.Email, user.Password).Login()
	})

	t.Run("Logout", func(t *testing.T) {
		client := suite.NewClient(t)
		loggedInUser := client.NewUser(user.Email, user.Password).Login()
		loggedInUser.Logout()

		// Try using revoked session
		loggedInUser.AuthedClient().POST("/auth/logout").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SessionUnauthorized).
			HasMessage("session not found or revoked")
	})
}
