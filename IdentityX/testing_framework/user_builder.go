package testing

import (
	"net/http"
	"testing"
)

// ============================================================================
// USER BUILDER - Declarative user management
// ============================================================================

type User struct {
	Email     string
	Password  string
	ProjectID string
	client    *Client
	auth      *AuthContext
	t         *testing.T
}

func (c *Client) User(email, password string) *User {
	return &User{
		Email:    email,
		Password: password,
		client:   c,
		t:        c.t,
	}
}

func (u *User) Register() *User {
	u.t.Helper()
	u.client.POST("/auth/register").
		WithBody(map[string]string{
			"email":    u.Email,
			"password": u.Password,
		}).
		Expect(http.StatusCreated).
		Success("go-auth-test", "Registered user")
	return u
}

func (u *User) ProjectRegister(projectID string) *User {
	u.t.Helper()
	u.client.POST("/projects/"+projectID+"/register").
		WithBody(map[string]interface{}{
			"email":         u.Email,
			"password":      u.Password,
			"custom_fields": []interface{}{},
		}).
		Expect(http.StatusCreated).
		Success("go-auth-test", "Registered user")
	return u
}

func (u *User) Login() *User {
	u.t.Helper()
	resp := u.client.POST("/auth/login").
		WithBody(map[string]string{
			"email":    u.Email,
			"password": u.Password,
		}).
		Expect(http.StatusOK)

	resp.Success("go-auth-test", "Logged in")
	u.auth = resp.Cookies()
	return u
}

func (u *User) ProjectLogin(projectID string) *User {
	u.t.Helper()
	resp := u.client.POST("/projects/" + projectID + "/login").
		WithBody(map[string]string{
			"email":    u.Email,
			"password": u.Password,
		}).
		Expect(http.StatusOK)

	resp.Success("go-auth-test", "Logged in")
	u.auth = resp.Cookies()
	return u
}

func (u *User) Logout() *User {
	u.t.Helper()
	u.authedClient().POST("/auth/logout").
		Expect(http.StatusOK).
		Success("go-auth-test", "Logged out")
	return u
}

func (u *User) Refresh() *User {
	u.t.Helper()

	req := u.client.expect.POST("/auth/refresh").
		WithCookie("refresh_token", u.auth.RefreshToken)

	resp := req.Expect().Status(http.StatusOK)

	access := resp.Cookie("access_token")
	refresh := resp.Cookie("refresh_token")

	if access.Raw() == nil || refresh.Raw() == nil {
		u.t.Fatal("Expected auth cookies after refresh but got nil")
		return u
	}

	u.auth = &AuthContext{
		AccessToken:  access.Value().Raw(),
		RefreshToken: refresh.Value().Raw(),
	}
	return u
}

func (u *User) authedClient() *Client {
	return u.client.Auth(u.auth)
}

func (u *User) AuthedClient() *Client {
	return u.authedClient()
}

func (u *User) CreateProject(name string) *User {
	u.t.Helper()
	resp := u.authedClient().POST("/projects").
		WithBody(map[string]interface{}{
			"project_name": name,
			"metadata":     map[string]string{"env": "test"},
		}).
		Expect(http.StatusCreated)

	resp.Success("go-auth-test", "Created project")
	u.ProjectID = resp.Data().Value("id").String().Raw()
	return u
}
