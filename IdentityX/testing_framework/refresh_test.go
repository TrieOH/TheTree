package testing

import (
	"net/http"
	"testing"
)

func testRefresh(t *testing.T, suite *TestSuite) {
	client := suite.Client(t)
	user := client.User("refresh@mail.com", ValidPassword).Register().Login()

	oldAccess := user.auth.AccessToken
	oldRefresh := user.auth.RefreshToken

	user.Refresh()

	if oldAccess == user.auth.AccessToken {
		t.Error("Access token should change after refresh")
	}

	if oldRefresh == user.auth.RefreshToken {
		t.Error("Refresh token should change after refresh")
	}

	// Old tokens should be invalid
	oldClient := client.Auth(&AuthContext{
		AccessToken:  oldAccess,
		RefreshToken: oldRefresh,
	})

	oldClient.GET("/sessions").
		Expect(http.StatusUnauthorized)
}
