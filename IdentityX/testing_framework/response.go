package testing

import (
	"GoAuth/internal/apierr"
	"strings"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
)

// ============================================================================
// RESPONSE - Chainable response assertions
// ============================================================================

type Response struct {
	resp   *httpexpect.Response
	t      *testing.T
	status int
}

func (r *Response) ValidationError(expectedErrors ...string) *Response {
	r.t.Helper()

	obj := r.resp.JSON().Object()

	// STRUCTURAL: fail fast
	obj.Value("module").String().IsEqual("validation")
	obj.Value("message").String().IsEqual("Validation failed")

	trace := r.RequireTrace()
	trace.Length().IsEqual(len(expectedErrors))

	// CONTENT: non-fatal diagnostics
	for i, err := range expectedErrors {
		actual := trace.Value(i).String().Raw()
		assert.Contains(r.t, actual, err, "trace[%d] mismatch", i)
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

func (r *Response) JSON() *httpexpect.Object {
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
			assert.Fail(r.t, "missing trace entry", "expected trace to contain %q, but it did not.\ntrace=%v", exp, raw)
		}
	}

	return r
}

func (r *Response) HasMessage(expected string) *Response {
	r.t.Helper()
	msg := r.resp.JSON().Object().Value("message").String().Raw()
	assert.Contains(r.t, msg, expected, "expected message to contain %q, but got %q", expected, msg)
	return r
}

func (r *Response) HasModule(expected string) *Response {
	r.t.Helper()
	msg := r.resp.JSON().Object().Value("module").String().Raw()
	assert.Contains(r.t, msg, expected, "expected module to contain %q, but got %q", expected, msg)
	return r
}

func (r *Response) HasErrID(expected apierr.ID) *Response {
	r.t.Helper()
	errID := r.resp.JSON().Object().Value("error_id").String().Raw()
	assert.Equal(r.t, string(expected), errID, "expected error id %q, but it was %q", string(expected), errID)
	return r
}
