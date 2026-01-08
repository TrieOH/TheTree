package testing

import (
	"GoAuth/internal/apierr"
	"strings"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// ============================================================================
// RESPONSE - Chainable response assertions
// ============================================================================

type Response struct {
	resp   *httpexpect.Response
	t      *testing.T
	status int
}

func (r *Response) Success(module, message string) *Response {
	r.t.Helper()
	obj := r.resp.JSON().Object()
	obj.Value("module").String().IsEqual(module)
	obj.Value("message").String().IsEqual(message)
	obj.Value("code").Number().IsEqual(r.status)
	return r
}

func (r *Response) Error(module, message string) *Response {
	r.t.Helper()
	return r.Success(module, message)
}

func (r *Response) ValidationError(expectedErrors ...string) *Response {
	r.t.Helper()
	obj := r.resp.JSON().Object()
	obj.Value("module").String().IsEqual("validation")
	obj.Value("message").String().IsEqual("Validation failed")

	trace := obj.Value("trace").Array()
	trace.Length().IsEqual(len(expectedErrors))

	for i, err := range expectedErrors {
		trace.Value(i).String().Contains(err)
	}

	return r
}

func (r *Response) Data() *httpexpect.Object {
	r.t.Helper()
	return r.resp.JSON().Object().Value("data").Object()
}

func (r *Response) Value() *httpexpect.Value {
	r.t.Helper()
	return r.resp.JSON().Object().Value("data")
}

func (r *Response) DataArray() *httpexpect.Array {
	r.t.Helper()
	return r.resp.JSON().Object().Value("data").Array()
}

func (r *Response) Cookies() *AuthContext {
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

func (r *Response) Trace() *httpexpect.Array {
	r.t.Helper()
	return r.resp.JSON().Object().Value("trace").Array()
}

func (r *Response) TraceContains(expected ...string) *Response {
	r.t.Helper()

	trace := r.Trace()
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
			r.t.Fatalf("expected trace to contain %q, but it did not.\ntrace=%v", exp, raw)
		}
	}

	return r
}

func (r *Response) MessageContains(expected string) *Response {
	r.t.Helper()

	msg := r.resp.JSON().Object().Value("message").String().Raw()

	if !strings.Contains(msg, expected) {
		r.t.Fatalf("expected message to contain %q, but it did not.\nmessage=%v", expected, msg)
	}

	return r
}

func (r *Response) ExpectErrorID(expected apierr.ID) *Response {
	r.t.Helper()

	errID := r.resp.JSON().Object().Value("error_id").String().Raw()

	if errID != string(expected) {
		r.t.Fatalf("expected error id %q, but it was %q", expected, errID)
	}

	return r
}
