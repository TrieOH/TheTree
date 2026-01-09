package testing

import (
	"net/http"
	"testing"
)

func testLogin(t *testing.T, suite *TestSuite) {
	// Create user in parent test context
	client := suite.Client(t)
	user := client.User("login@mail.com", ValidPassword).Register()

	t.Run("WrongPassword", func(t *testing.T) {
		client := suite.Client(t)
		client.POST("/auth/login").
			WithBody(map[string]string{
				"email":    user.Email,
				"password": "WrongPass123!",
			}).
			Expect(http.StatusUnauthorized).
			Error("go-auth-test", "invalid email or password")
	})

	t.Run("WrongEmail", func(t *testing.T) {
		client := suite.Client(t)
		client.POST("/auth/login").
			WithBody(map[string]string{
				"email":    "wrong@mail.com",
				"password": user.Password,
			}).
			Expect(http.StatusUnauthorized).
			Error("go-auth-test", "invalid email or password")
	})

	t.Run("Success", func(t *testing.T) {
		client := suite.Client(t)
		client.User(user.Email, user.Password).Login()
	})

	t.Run("Logout", func(t *testing.T) {
		client := suite.Client(t)
		loggedInUser := client.User(user.Email, user.Password).Login()
		loggedInUser.Logout()

		// Try using revoked session
		loggedInUser.AuthedClient().POST("/auth/logout").
			Expect(http.StatusUnauthorized).
			Error("AuthMW", "session not found or revoked")
	})
}
