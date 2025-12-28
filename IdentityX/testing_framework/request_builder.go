package testing

import (
	"testing"

	"github.com/gavv/httpexpect/v2"
)

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

func (rb *RequestBuilder) Expect(status int) *Response {
	rb.t.Helper()
	return &Response{
		resp:   rb.req.Expect().Status(status),
		t:      rb.t,
		status: status,
	}
}

// ExpectStatus is an alias for Expect for clarity
func (rb *RequestBuilder) ExpectStatus(status int) *Response {
	return rb.Expect(status)
}
