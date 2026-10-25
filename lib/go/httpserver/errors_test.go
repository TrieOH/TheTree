package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"lib/database"

	"github.com/MintzyG/fun"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var (
	errAddSelf      = errors.New("users can't add themselves to projects")
	errBadCreds     = errors.New("invalid email or password")
	errTermsMissing = errors.New("you must accept the terms of service to register")
)

type outOfStockError struct{ items int }

func (e *outOfStockError) Error() string { return fmt.Sprintf("%d items are out of stock", e.items) }

// useTestErrorMap installs m and restores an empty map after the test.
func useTestErrorMap(t *testing.T, m ErrorMap) {
	t.Helper()
	UseErrorMap(m)
	t.Cleanup(func() {
		fun.ResetAppErrorMapper()
		errorMapMu.Lock()
		errorMap = nil
		errorMapMu.Unlock()
	})
}

func TestErrorMapResolve(t *testing.T) {
	m := append(ErrorMap{
		Is(errAddSelf, http.StatusBadRequest),
		Is(errBadCreds, http.StatusUnauthorized),
		Is(errTermsMissing, http.StatusBadRequest).Validation(),
		As[*outOfStockError](http.StatusConflict),
	}, sharedErrors...)

	cases := []struct {
		name   string
		err    error
		status int
		code   fun.ErrorCode
		msg    string
	}{
		{"sentinel", errAddSelf, 400, fun.CodeBadRequest, errAddSelf.Error()},
		{"wrapped: context never reaches the client",
			fmt.Errorf("add member: %w", errAddSelf), 400, fun.CodeBadRequest, errAddSelf.Error()},
		{"validation code", errTermsMissing, 400, fun.CodeValidation, errTermsMissing.Error()},
		{"typed", fmt.Errorf("checkout: %w", &outOfStockError{items: 2}), 409, fun.CodeConflict, "2 items are out of stock"},
		{"most specific wins over its class",
			fmt.Errorf("login: %w", fmt.Errorf("%w: %w", errBadCreds, database.ErrNotFound)), 401, fun.CodeUnauthorized, errBadCreds.Error()},
		{"class alone", fmt.Errorf("get: %w", database.ErrNotFound), 404, fun.CodeNotFound, database.ErrNotFound.Error()},
		{"unregistered constraint maps by class, hides the constraint",
			&database.ConstraintError{Class: database.ErrUniqueViolation, Constraint: "uniq_secret"}, 409, fun.CodeConflict, database.ErrUniqueViolation.Error()},
		{"unmapped is a generic 500", errors.New("pq: connection reset"), 500, fun.CodeInternal, "internal error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, code, msg := m.resolve(tc.err)
			if status != tc.status || code != tc.code || msg != tc.msg {
				t.Fatalf("got (%d, %s, %q), want (%d, %s, %q)", status, code, msg, tc.status, tc.code, tc.msg)
			}
		})
	}
}

// TestWriteErrorKeepsEnvelope — a named error goes out in the fun envelope
// with the mapped status, code and message.
func TestWriteErrorKeepsEnvelope(t *testing.T) {
	useTestErrorMap(t, ErrorMap{Is(errAddSelf, http.StatusBadRequest)})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/projects/x/members", nil)
	WriteError(rec, req, fmt.Errorf("add member: %w", errAddSelf))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
	var body struct {
		Code  int `json:"code"`
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Code != 400 || body.Error.Code != "BAD_REQUEST" || body.Error.Message != errAddSelf.Error() {
		t.Fatalf("unexpected envelope: %s", rec.Body.String())
	}
}

// TestTraceMiddlewareLeavesClientErrorsUnset — a 4xx is recorded on the
// span but does not mark it failed.
func TestTraceMiddlewareLeavesClientErrorsUnset(t *testing.T) {
	useTestErrorMap(t, ErrorMap{Is(errAddSelf, http.StatusBadRequest)})
	rec := tracetest.NewSpanRecorder()
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec)))
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	handler := TraceMiddleware()(func(context.Context, http.ResponseWriter, *http.Request, any) (any, error) {
		return nil, errAddSelf
	}, "AddProjectMember")
	req := httptest.NewRequest(http.MethodPost, "/projects/x/members", nil)
	_, _ = handler(req.Context(), httptest.NewRecorder(), req, nil)

	spans := rec.Ended()
	if len(spans) != 1 {
		t.Fatalf("want 1 span, got %d", len(spans))
	}
	if spans[0].Status().Code == codes.Error {
		t.Fatal("a 4xx marked the span as failed")
	}
	if len(spans[0].Events()) == 0 {
		t.Fatal("the error was not recorded on the span")
	}
}
