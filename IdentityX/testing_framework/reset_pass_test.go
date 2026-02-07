package testing

import (
	"GoAuth/internal/apierr"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"testing"
	"time"
)

func testResetPassword(t *testing.T, suite *TestSuite) {
	baseClient := suite.NewClient(t)

	// 1. Global User Reset Password
	t.Run("GlobalUserResetPassword", func(t *testing.T) {
		client := baseClient.WithT(t)
		userEmail := "global-reset@mail.com"
		newPassword := "NewSecurePassword123!"

		// Register and Login to have an active session
		client.WithCredentials(userEmail, ValidPassword).Register()
		userSession := client.WithCredentials(userEmail, ValidPassword).Login()

		// Request forgot password
		client.POST("/auth/forgot-password").
			WithBody(map[string]string{
				"email": userEmail,
			}).
			Expect(http.StatusOK).
			HasMessage("forgot password email sent")

		time.Sleep(time.Millisecond * 500)

		link, err := GetLatestResetPasswordLink(userEmail)
		if err != nil {
			t.Fatalf("Failed getting reset password link: %v", err.Error())
		}

		u, err := url.Parse(link)
		if err != nil {
			t.Fatalf("Failed parsing reset password link: %v", err)
		}
		token := u.Query().Get("token")

		// Reset password
		client.POST("/auth/reset-password").
			WithQuery("token", token).
			WithBody(map[string]string{
				"new_password": newPassword,
			}).
			Expect(http.StatusOK).
			HasMessage("password reset successfully")

		// Check if old session is invalidated
		userSession.GET("/sessions/me").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SessionRevoked)

		// Try to login with old password - should fail
		client.POST("/auth/login").
			WithBody(map[string]string{
				"email":    userEmail,
				"password": ValidPassword,
			}).
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.AuthInvalidCredentials)

		// Login with new password - should succeed
		client.WithCredentials(userEmail, newPassword).
			Login()
	})

	// 2. Project User Reset Password
	t.Run("ProjectUserResetPassword", func(t *testing.T) {
		client := baseClient.WithT(t)
		projectOwnerEmail := "project-owner-reset@mail.com"
		projectOwner := suite.NewClient(t).WithCredentials(projectOwnerEmail, ValidPassword).Register().Login()
		project := projectOwner.CreateProject("ResetProject")
		projectID := project.projectID

		projectUserEmail := "project-user-reset@mail.com"
		projectUserPassword := ValidPassword
		projectUserNewPassword := "ProjectNewPass123!"

		// Register and login project user
		projectOwner.WithCredentials(projectUserEmail, projectUserPassword).ProjectRegister(projectID)
		projectUserSession := suite.NewClient(t).WithCredentials(projectUserEmail, projectUserPassword).ProjectLogin(projectID)

		// Request forgot password for project user
		client.POST("/auth/forgot-password").
			WithBody(map[string]interface{}{
				"email":      projectUserEmail,
				"project_id": projectID,
			}).
			Expect(http.StatusOK).
			HasMessage("forgot password email sent")

		time.Sleep(time.Millisecond * 500)

		link, err := GetLatestResetPasswordLink(projectUserEmail)
		if err != nil {
			t.Fatalf("Failed getting reset password link: %v", err.Error())
		}

		u, err := url.Parse(link)
		if err != nil {
			t.Fatalf("Failed parsing reset password link: %v", err)
		}
		token := u.Query().Get("token")

		// Reset password
		client.POST("/auth/reset-password").
			WithQuery("token", token).
			WithBody(map[string]string{
				"new_password": projectUserNewPassword,
			}).
			Expect(http.StatusOK).
			HasMessage("password reset successfully")

		// Check if project user session is invalidated
		projectUserSession.GET("/sessions/me").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.SessionRevoked)

		// Try to login with old password - should fail
		client.POST("/projects/" + projectID + "/login").
			WithBody(map[string]string{
				"email":    projectUserEmail,
				"password": projectUserPassword,
			}).
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.AuthInvalidCredentials)

		// Login with new password - should succeed
		client.WithCredentials(projectUserEmail, projectUserNewPassword).
			ProjectLogin(projectID)
	})

	// 3. Token Reuse Prevention
	t.Run("TokenReusePrevention", func(t *testing.T) {
		client := baseClient.WithT(t)
		email := "reuse-token@mail.com"
		client.WithCredentials(email, ValidPassword).Register()

		client.POST("/auth/forgot-password").
			WithBody(map[string]string{
				"email": email,
			}).
			Expect(http.StatusOK)

		time.Sleep(time.Millisecond * 500)
		link, err := GetLatestResetPasswordLink(email)
		if err != nil {
			t.Fatalf("Failed getting reset password link: %v", err.Error())
		}
		u, _ := url.Parse(link)
		token := u.Query().Get("token")

		// First use
		client.POST("/auth/reset-password").
			WithQuery("token", token).
			WithBody(map[string]string{
				"new_password": "NewPassword1!",
			}).
			Expect(http.StatusOK)

		// Second use
		client.POST("/auth/reset-password").
			WithQuery("token", token).
			WithBody(map[string]string{
				"new_password": "NewPassword2!",
			}).
			Expect(http.StatusForbidden).
			HasErrID(apierr.AuthTokenAlreadyUsed)
	})
}

// MailHog types for decoding
type MHMessage struct {
	Content struct {
		Body string `json:"Body"`
	} `json:"Content"`
	Raw struct {
		To []string `json:"To"`
	} `json:"Raw"`
}

type MHResponse struct {
	Items []MHMessage `json:"items"`
}

func GetLatestResetPasswordLink(toEmail string) (string, error) {
	resp, err := http.Get("http://mailhog:8025/api/v2/messages")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var mh MHResponse
	if err := json.NewDecoder(resp.Body).Decode(&mh); err != nil {
		return "", err
	}

	if len(mh.Items) == 0 {
		return "", errors.New("no emails found")
	}

	// The link in PasswordReset is /reset?token=...
	re := regexp.MustCompile(`href="([^"]+)"`)

	// We want the latest email for this user
	for _, item := range mh.Items {
		// Check if it's for the right email
		isForEmail := false
		for _, to := range item.Raw.To {
			if to == toEmail {
				isForEmail = true
				break
			}
		}
		if !isForEmail {
			continue
		}

		match := re.FindStringSubmatch(item.Content.Body)
		if len(match) >= 2 && (regexp.MustCompile(`/reset\?token=`).MatchString(match[1])) {
			return match[1], nil
		}
	}

	return "", errors.New("reset password link not found for " + toEmail)
}
