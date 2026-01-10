package testing

import (
	"GoAuth/internal/apierr"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/goforj/godump"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// CLIENT - HTTP client with optional authentication and credentials
// ============================================================================

type Client struct {
	expect *httpexpect.Expect
	t      *testing.T
	auth   *AuthContext

	// Credentials for re-authentication
	email     string
	password  string
	projectID string
}

type AuthContext struct {
	AccessToken  string
	RefreshToken string
}

// ----------------
// Factory methods
// ----------------

// WithCredentials creates a new client with credentials (unauthenticated initially)
func (c *Client) WithCredentials(email, password string) *Client {
	return &Client{
		expect:    c.expect,
		t:         c.t,
		email:     email,
		password:  password,
		projectID: c.projectID,
	}
}

// WithAuth creates a new client with the given auth context
func (c *Client) WithAuth(auth *AuthContext) *Client {
	return &Client{
		expect:    c.expect,
		t:         c.t,
		auth:      auth,
		email:     c.email,
		password:  c.password,
		projectID: c.projectID,
	}
}

// WithT creates a new client with a different testing.T (for subtests)
func (c *Client) WithT(t *testing.T) *Client {
	return &Client{
		expect:    c.expect,
		t:         t,
		auth:      c.auth,
		email:     c.email,
		password:  c.password,
		projectID: c.projectID,
	}
}

// ----------------
// Auth operations
// ----------------

// Register registers the client's credentials and returns the same client
func (c *Client) Register() *Client {
	c.t.Helper()
	c.POST("/auth/register").
		WithBody(map[string]string{
			"email":    c.email,
			"password": c.password,
		}).
		Expect(http.StatusCreated).
		HasModule("go-auth-test").
		HasMessage("Registered user")
	return c
}

// ProjectRegister registers the client in a specific project
func (c *Client) ProjectRegister(projectID string) *Client {
	c.t.Helper()
	c.POST("/projects/" + projectID + "/register").
		WithBody(map[string]interface{}{
			"email":    c.email,
			"password": c.password,
		}).
		Expect(http.StatusCreated).
		HasModule("go-auth-test").
		HasMessage("Registered user")
	return c
}

// Login authenticates and returns a new client with auth cookies
func (c *Client) Login() *Client {
	c.t.Helper()
	auth := c.POST("/auth/login").
		WithBody(map[string]string{
			"email":    c.email,
			"password": c.password,
		}).
		Expect(http.StatusOK).
		HasModule("go-auth-test").
		HasMessage("Logged in").
		AuthCookies()

	return c.WithAuth(auth)
}

// ProjectLogin authenticates in a specific project
func (c *Client) ProjectLogin(projectID string) *Client {
	c.t.Helper()
	auth := c.POST("/projects/" + projectID + "/login").
		WithBody(map[string]string{
			"email":    c.email,
			"password": c.password,
		}).
		Expect(http.StatusOK).
		HasModule("go-auth-test").
		HasMessage("Logged in").
		AuthCookies()

	return c.WithAuth(auth)
}

// Logout logs out the current session
func (c *Client) Logout() *Client {
	c.t.Helper()
	c.POST("/auth/logout").
		Expect(http.StatusOK).
		HasModule("go-auth-test").
		HasMessage("Logged out")
	return c
}

// Refresh refreshes the auth tokens and returns a new client with updated auth
func (c *Client) Refresh() *Client {
	c.t.Helper()

	if c.auth == nil {
		c.t.Fatal("Refresh called on unauthenticated client")
		return c
	}

	req := c.expect.POST("/auth/refresh").
		WithCookie("refresh_token", c.auth.RefreshToken)

	resp := req.Expect().Status(http.StatusOK)

	access := resp.Cookie("access_token")
	refresh := resp.Cookie("refresh_token")

	if access.Raw() == nil || refresh.Raw() == nil {
		c.t.Fatal("Expected auth cookies after refresh but got nil")
		return c
	}

	return c.WithAuth(&AuthContext{
		AccessToken:  access.Value().Raw(),
		RefreshToken: refresh.Value().Raw(),
	})
}

// ----------------
// Project operations
// ----------------

// CreateProject creates a project and stores the ID for chaining
func (c *Client) CreateProject(name string) *Client {
	c.t.Helper()
	resp := c.POST("/projects").
		WithBody(map[string]interface{}{
			"project_name": name,
			"metadata":     map[string]string{"env": "test"},
		}).
		Expect(http.StatusCreated).
		HasMessage("Created project")

	projectID := resp.RequireDataObject().Value("id").String().NotEmpty().Raw()

	return &Client{
		expect:    c.expect,
		t:         c.t,
		auth:      c.auth,
		email:     c.email,
		password:  c.password,
		projectID: projectID,
	}
}

// ProjectID returns the current project ID (useful for assertions)
func (c *Client) ProjectID() string {
	return c.projectID
}

// ----------------
// Request builders
// ----------------

func (c *Client) POST(path string) *RequestBuilder {
	return c.newRequest("POST", path)
}

func (c *Client) GET(path string) *RequestBuilder {
	return c.newRequest("GET", path)
}

func (c *Client) PATCH(path string) *RequestBuilder {
	return c.newRequest("PATCH", path)
}

func (c *Client) DELETE(path string) *RequestBuilder {
	return c.newRequest("DELETE", path)
}

func (c *Client) newRequest(method, path string) *RequestBuilder {
	req := c.expect.Request(method, path).
		WithHeader("Content-Type", "application/json")

	if c.auth != nil {
		req = req.
			WithCookie("access_token", c.auth.AccessToken).
			WithCookie("refresh_token", c.auth.RefreshToken)
	}

	return &RequestBuilder{
		req: req,
		t:   c.t,
	}
}

// ============================================================================
// REQUEST BUILDER - Fluent request construction
// ============================================================================

type RequestBuilder struct {
	req *httpexpect.Request
	t   *testing.T
}

func (rb *RequestBuilder) WithBody(body interface{}) *RequestBuilder {
	rb.req = rb.req.WithJSON(body)
	return rb
}

func (rb *RequestBuilder) WithQuery(key, value string) *RequestBuilder {
	rb.req = rb.req.WithQuery(key, value)
	return rb
}

func (rb *RequestBuilder) WithCookie(key, value string) *RequestBuilder {
	rb.req = rb.req.WithCookie(key, value)
	return rb
}

func (rb *RequestBuilder) Expect(status int) *Response {
	rb.t.Helper()

	httpResp := rb.req.Expect()
	actualStatus := httpResp.Raw().StatusCode

	// Capture body early before any other method consumes it
	bodyStr := httpResp.Body().Raw()

	r := &Response{
		resp:       httpResp,
		t:          rb.t,
		status:     status,
		bodyCache:  bodyStr,
		failed:     false,
		dumpedOnce: false,
	}

	// Check if status matches, if not dump immediately
	httpResp.Status(status)
	if actualStatus != status {
		r.failed = true
		r.dumpOnce()
	}

	return r
}

// ExpectStatus is an alias for Expect for clarity
func (rb *RequestBuilder) ExpectStatus(status int) *Response {
	return rb.Expect(status)
}

// ============================================================================
// RESPONSE - Chainable response assertions
// ============================================================================

type Response struct {
	resp       *httpexpect.Response
	t          *testing.T
	status     int
	bodyCache  string
	failed     bool
	dumpedOnce bool
}

func (r *Response) dumpOnce() {
	if r.dumpedOnce {
		return
	}
	r.dumpedOnce = true

	r.t.Logf("\n=== HTTP Response Dump (status %d) ===", r.status)

	// Try to parse as JSON for pretty printing
	var parsedBody interface{}
	if err := json.Unmarshal([]byte(r.bodyCache), &parsedBody); err == nil {
		godump.Dump(parsedBody)
	} else {
		// Not JSON, dump as string
		godump.Dump(r.bodyCache)
	}
}

func (r *Response) ValidationError(expectedErrors ...string) *Response {
	r.t.Helper()

	obj := r.resp.JSON().Object()

	// STRUCTURAL: fail fast
	obj.Value("module").String().IsEqual("validation")
	obj.Value("message").String().IsEqual("Validation failed")

	trace := r.RequireTrace()
	trace.Length().IsEqual(len(expectedErrors))

	// CONTENT: soft assertions
	for i, err := range expectedErrors {
		actual := trace.Value(i).String().Raw()
		if !assert.Contains(r.t, actual, err, "trace[%d] mismatch", i) {
			r.failed = true
			r.dumpOnce()
		}
	}

	return r
}

func (r *Response) RequireDataArray() *httpexpect.Array {
	r.t.Helper()
	return r.resp.JSON().Object().Value("data").Array()
}

func (r *Response) RequireDataObject() *httpexpect.Object {
	r.t.Helper()
	return r.resp.JSON().Object().Value("data").Object()
}

func (r *Response) RequireDataValue() *httpexpect.Value {
	r.t.Helper()
	return r.resp.JSON().Object().Value("data")
}

func (r *Response) AuthCookies() *AuthContext {
	r.t.Helper()
	access := r.resp.Cookie("access_token")
	refresh := r.resp.Cookie("refresh_token")

	if access.Raw() == nil || refresh.Raw() == nil {
		r.t.Fatal("Expected auth cookies but got nil")
		return nil
	}

	return &AuthContext{
		AccessToken:  access.Value().Raw(),
		RefreshToken: refresh.Value().Raw(),
	}
}

func (r *Response) Resp() *httpexpect.Response {
	r.t.Helper()
	return r.resp
}

func (r *Response) JSON() *httpexpect.Value {
	r.t.Helper()
	return r.resp.JSON()
}

func (r *Response) JSONObj() *httpexpect.Object {
	r.t.Helper()
	return r.resp.JSON().Object()
}

func (r *Response) RequireTrace() *httpexpect.Array {
	r.t.Helper()
	return r.resp.JSON().Object().Value("trace").Array()
}

func (r *Response) TraceContains(expected ...string) *Response {
	r.t.Helper()

	trace := r.RequireTrace()
	raw := trace.Raw()

	for _, exp := range expected {
		found := false

		for _, v := range raw {
			s, ok := v.(string)
			if ok && strings.Contains(s, exp) {
				found = true
				break
			}
		}

		if !found {
			r.t.Errorf("missing trace entry: expected trace to contain %q, but it did not.\ntrace=%v", exp, raw)
			r.failed = true
			r.dumpOnce()
		}
	}

	return r
}

func (r *Response) HasMessage(expected string) *Response {
	r.t.Helper()
	msg := r.resp.JSON().Object().Value("message").String().Raw()
	if !assert.Contains(r.t, msg, expected, "expected message to contain %q, but got %q", expected, msg) {
		r.failed = true
		r.dumpOnce()
	}
	return r
}

func (r *Response) HasModule(expected string) *Response {
	r.t.Helper()
	msg := r.resp.JSON().Object().Value("module").String().Raw()
	if !assert.Contains(r.t, msg, expected, "expected module to contain %q, but got %q", expected, msg) {
		r.failed = true
		r.dumpOnce()
	}
	return r
}

func (r *Response) HasErrID(expected apierr.ID) *Response {
	r.t.Helper()
	errID := r.resp.JSON().Object().Value("error_id").String().Raw()
	if !assert.Equal(r.t, string(expected), errID, "expected error id %q, but it was %q", string(expected), errID) {
		r.failed = true
		r.dumpOnce()
	}
	return r
}
