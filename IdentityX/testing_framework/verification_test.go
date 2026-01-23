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

	"github.com/stretchr/testify/assert"
)

func testVerification(t *testing.T, suite *TestSuite) {
	client := suite.NewClient(t)
	user := client.WithCredentials("verification@mail.com", ValidPassword).Register().Login().CreateProject("verification")

	time.Sleep(time.Millisecond * 100)

	link, err := GetLatestVerificationLink()
	if err != nil {
		t.Fatalf("Failed getting verification link: %v", err.Error())
	}

	if link == "" {
		t.Fatalf("Verification link is empty")
	}

	u, err := url.Parse(link)
	if err != nil {
		t.Fatalf("Failed parsing verification link: %v", err)
	}

	t.Run("ResendVerificationEmail", func(t *testing.T) {
		user.WithT(t).POST("/auth/verify/resend").
			Expect(http.StatusOK).
			HasMessage("verification email resent successfully")

		time.Sleep(time.Millisecond * 100)

		link2, err2 := GetLatestVerificationLink()
		if err2 != nil {
			t.Fatalf("Failed getting verification link: %v", err2.Error())
		}

		if link2 == "" {
			t.Fatalf("Verification link is empty")
		}

		u2, err2 := url.Parse(link2)
		if err2 != nil {
			t.Fatalf("Failed parsing verification link: %v", err2)
		}

		assert.NotEqualf(t, u2.Query().Get("token"), u.Query().Get("token"), "tokens shouldn't be equal")
	})

	t.Run("VerifyUser", func(t *testing.T) {
		user.WithT(t).POST(u.Path).
			WithQuery("token", u.Query().Get("token")).
			Expect(http.StatusOK).
			HasMessage("user verified, please refresh")
	})

	t.Run("VerifyUserAgain", func(t *testing.T) {
		user.WithT(t).POST(u.Path).
			WithQuery("token", u.Query().Get("token")).
			Expect(http.StatusOK).
			HasMessage("user verified, please refresh")
	})

	t.Run("ResendVerificationEmailNotAllowed", func(t *testing.T) {
		user.WithT(t).POST("/auth/verify/resend").
			Expect(http.StatusForbidden).
			HasErrID(apierr.AuthAlreadyVerified).
			HasMessage("user already verified")
	})

	t.Run("SessionInfoBeforeRefreshed", func(t *testing.T) {
		user := suite.NewClient(t).WithAuth(user.auth)

		data := user.GET("/sessions/me").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"refresh_expire_date": AnyNumber{},
			"access": map[string]interface{}{
				"iss": "GoAuth",
				"exp": AnyNumber{},
				"iat": AnyNumber{},
				"jti": AnyUUID{},
				"sub": map[string]interface{}{
					"id":          AnyUUID{},
					"email":       "verification@mail.com",
					"project_id":  nil,
					"user_type":   "client",
					"metadata":    nil,
					"session_id":  AnyUUID{},
					"user_agent":  AnyString{},
					"user_ip":     AnyString{},
					"is_verified": false,
					"verified_at": nil,
				},
			},
		}

		Validate(t, data, spec)
	})

	t.Run("SessionInfoRefreshed", func(t *testing.T) {
		user := suite.NewClient(t).WithCredentials("verification@mail.com", ValidPassword).Login()

		data := user.GET("/sessions/me").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"refresh_expire_date": AnyNumber{},
			"access": map[string]interface{}{
				"iss": "GoAuth",
				"exp": AnyNumber{},
				"iat": AnyNumber{},
				"jti": AnyUUID{},
				"sub": map[string]interface{}{
					"id":          AnyUUID{},
					"email":       "verification@mail.com",
					"project_id":  nil,
					"user_type":   "client",
					"metadata":    nil,
					"session_id":  AnyUUID{},
					"user_agent":  AnyString{},
					"user_ip":     AnyString{},
					"is_verified": true,
					"verified_at": AnyDate{},
				},
			},
		}

		Validate(t, data, spec)
	})

	var projectUser *Client
	projectUser = client.WithCredentials("verification-project@mail.com", ValidPassword).ProjectRegister(user.projectID).ProjectLogin(user.projectID)

	time.Sleep(time.Millisecond * 100)

	link2, err := GetLatestVerificationLink()
	if err != nil {
		t.Fatalf("Failed getting verification link: %v", err.Error())
	}

	if link2 == "" {
		t.Fatalf("Verification link is empty")
	}

	u2, err := url.Parse(link2)
	if err != nil {
		t.Fatalf("Failed parsing verification link: %v", err)
	}

	t.Run("ResendVerificationEmailProjectUser", func(t *testing.T) {
		projectUser.WithT(t).POST("/auth/verify/resend").
			Expect(http.StatusOK).
			HasMessage("verification email resent successfully")

		time.Sleep(time.Millisecond * 100)

		link3, err2 := GetLatestVerificationLink()
		if err2 != nil {
			t.Fatalf("Failed getting verification link: %v", err2.Error())
		}

		if link3 == "" {
			t.Fatalf("Verification link is empty")
		}

		u3, err2 := url.Parse(link3)
		if err2 != nil {
			t.Fatalf("Failed parsing verification link: %v", err2)
		}

		assert.NotEqualf(t, u3.Query().Get("token"), u2.Query().Get("token"), "tokens shouldn't be equal")
	})

	t.Run("VerifyProjectUser", func(t *testing.T) {
		projectUser.WithT(t).POST(u2.Path).
			WithQuery("token", u2.Query().Get("token")).
			Expect(http.StatusOK).
			HasMessage("user verified, please refresh")
	})

	t.Run("VerifyProjectUserAgain", func(t *testing.T) {
		projectUser.WithT(t).POST(u2.Path).
			WithQuery("token", u2.Query().Get("token")).
			Expect(http.StatusOK).
			HasMessage("user verified, please refresh")
	})

	t.Run("ResendVerificationEmailNotAllowedProjectUser", func(t *testing.T) {
		projectUser.WithT(t).POST("/auth/verify/resend").
			Expect(http.StatusForbidden).
			HasErrID(apierr.AuthAlreadyVerified).
			HasMessage("user already verified")
	})

	t.Run("ProjUserSessionInfoBeforeRefreshed", func(t *testing.T) {
		data := projectUser.WithT(t).GET("/sessions/me").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"refresh_expire_date": AnyNumber{},
			"access": map[string]interface{}{
				"iss": "GoAuth",
				"exp": AnyNumber{},
				"iat": AnyNumber{},
				"jti": AnyUUID{},
				"sub": map[string]interface{}{
					"id":          AnyUUID{},
					"email":       "verification-project@mail.com",
					"project_id":  AsString{Value: user.projectID, Matcher: AnyUUID{}},
					"user_type":   "project",
					"metadata":    map[string]interface{}{},
					"session_id":  AnyUUID{},
					"user_agent":  AnyString{},
					"user_ip":     AnyString{},
					"is_verified": false,
					"verified_at": nil,
				},
			},
		}

		Validate(t, data, spec)
	})

	t.Run("ProjUserSessionInfoRefreshed", func(t *testing.T) {
		projectUser := suite.NewClient(t).WithCredentials("verification-project@mail.com", ValidPassword).ProjectLogin(user.projectID)

		data := projectUser.GET("/sessions/me").
			Expect(http.StatusOK).
			RequireDataValue()

		spec := map[string]interface{}{
			"refresh_expire_date": AnyNumber{},
			"access": map[string]interface{}{
				"iss": "GoAuth",
				"exp": AnyNumber{},
				"iat": AnyNumber{},
				"jti": AnyUUID{},
				"sub": map[string]interface{}{
					"id":          AnyUUID{},
					"email":       "verification-project@mail.com",
					"project_id":  AsString{Value: user.projectID, Matcher: AnyUUID{}},
					"user_type":   "project",
					"metadata":    map[string]interface{}{},
					"session_id":  AnyUUID{},
					"user_agent":  AnyString{},
					"user_ip":     AnyString{},
					"is_verified": true,
					"verified_at": AnyDate{},
				},
			},
		}

		Validate(t, data, spec)
	})
}

type MailHogResponse struct {
	Items []struct {
		Content struct {
			Body string `json:"Body"`
		} `json:"Content"`
	} `json:"items"`
}

func GetLatestVerificationLink() (string, error) {
	resp, err := http.Get("http://mailhog:8025/api/v2/messages")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var mh MailHogResponse
	if err := json.NewDecoder(resp.Body).Decode(&mh); err != nil {
		return "", err
	}

	if len(mh.Items) == 0 {
		return "", errors.New("no emails found")
	}

	re := regexp.MustCompile(`href="([^"]+)"`)
	match := re.FindStringSubmatch(mh.Items[0].Content.Body)
	if len(match) < 2 {
		return "", errors.New("verification link not found")
	}

	return match[1], nil
}
