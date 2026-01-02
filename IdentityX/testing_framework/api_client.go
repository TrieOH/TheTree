package testing

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// ============================================================================
// API CLIENT - Fluent interface for making requests
// ============================================================================

type Client struct {
	expect *httpexpect.Expect
	t      *testing.T
	auth   *AuthContext
}

type AuthContext struct {
	AccessToken  string
	RefreshToken string
}

// Auth adds authentication to subsequent requests
func (c *Client) Auth(ctx *AuthContext) *Client {
	return &Client{
		expect: c.expect,
		t:      c.t,
		auth:   ctx,
	}
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
