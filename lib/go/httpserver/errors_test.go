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
	"lib/fielderr"
	"lib/validator"

	idx "sdk/identityx"

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

type upstreamError struct{ status int }

func (e *upstreamError) Error() string { return fmt.Sprintf("upstream said %d", e.status) }

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
		AsStatus(func(e *upstreamError) int {
			if e.status >= 500 {
				return http.StatusBadGateway
			}
			return http.StatusBadRequest
		}),
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
		{"missing row reads as the resource's error",
			&database.NotFoundError{Named: errors.New("form not found")}, 404, fun.CodeNotFound, "form not found"},
		{"registered violation maps by its named error",
			&database.ViolationError{Named: errTermsMissing, Class: database.ErrCheckViolation}, 400, fun.CodeValidation, errTermsMissing.Error()},
		{"batch lists its items",
			&database.BatchError{Items: []database.BatchItemError{{ID: "f1", Err: errAddSelf}}}, 400, fun.CodeValidation, "f1: " + errAddSelf.Error()},
		{"validation lists its fields in the message",
			&validator.Error{Fields: []fielderr.Field{{Field: "name", Message: "this field is required"}}}, 400, fun.CodeValidation, "invalid request body: name: this field is required"},
		{"no identity is a 401", fmt.Errorf("handler: %w", idx.ErrNoIdentity), 401, fun.CodeUnauthorized, idx.ErrNoIdentity.Error()},
		{"status from the error's data", &upstreamError{status: 503}, 502, fun.CodeBadGateway, "upstream said 503"},
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
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/projects/x/members", nil)
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

// TestWriteMiddlewareErrorKeepsItsStatus — a middleware's rejection keeps
// the middleware's status whatever the map says, reads as the mapped
// link's message, and a 5xx never shows its cause.
func TestWriteMiddlewareErrorKeepsItsStatus(t *testing.T) {
	useTestErrorMap(t, ErrorMap{Is(errBadCreds, http.StatusForbidden)})

	cases := []struct {
		name   string
		status int
		err    error
		msg    string
	}{
		{"mapped error keeps the middleware's status", 401, fmt.Errorf("hook: %w", errBadCreds), errBadCreds.Error()},
		{"unmapped 4xx reads as its error", 401, errors.New("missing X-API-Key header"), "missing X-API-Key header"},
		{"5xx stays generic", 500, errors.New("panic: nil map"), "internal error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/projects", nil)
			writeMiddlewareError(rec, req, tc.status, tc.err)

			var body struct {
				Error struct {
					Message string `json:"message"`
				} `json:"error"`
			}
			err := json.Unmarshal(rec.Body.Bytes(), &body)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if rec.Code != tc.status || body.Error.Message != tc.msg {
				t.Fatalf("got (%d, %q), want (%d, %q)", rec.Code, body.Error.Message, tc.status, tc.msg)
			}
		})
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
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/projects/x/members", nil)
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

// TestWriteErrorListsFields — a fielderr.Error answers as its named error
// and lists its fields in the envelope.
func TestWriteErrorListsFields(t *testing.T) {
	useTestErrorMap(t, ErrorMap{Is(errAddSelf, http.StatusBadRequest)})

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/checkout", nil)
	WriteError(rec, req, fmt.Errorf("checkout: %w", fielderr.On(errAddSelf, "items[0].attendee", "required")))

	var body struct {
		Error struct {
			Message string `json:"message"`
			Fields  []struct {
				Field   string `json:"field"`
				Message string `json:"message"`
			} `json:"fields"`
		} `json:"error"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if rec.Code != http.StatusBadRequest || body.Error.Message != errAddSelf.Error() ||
		len(body.Error.Fields) != 1 || body.Error.Fields[0].Field != "items[0].attendee" {
		t.Fatalf("unexpected response %d: %s", rec.Code, rec.Body.String())
	}
}

// TestPlainErrorsDropTheEnvelope — with UsePlainErrors the body is the bare
// error object, with the mapped status, code, message and fields.
func TestPlainErrorsDropTheEnvelope(t *testing.T) {
	useTestErrorMap(t, ErrorMap{Is(errTermsMissing, http.StatusBadRequest).Validation()})
	UsePlainErrors()
	t.Cleanup(func() { plainErrors.Store(false) })

	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/auth/register", nil)
	WriteError(rec, req, fielderr.On(errTermsMissing, "accept_terms", "required"))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
	var body map[string]any
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	want := map[string]any{
		"code":    "VALIDATION_ERROR",
		"message": errTermsMissing.Error(),
		"fields":  []any{map[string]any{"field": "accept_terms", "message": "required"}},
	}
	if fmt.Sprint(body) != fmt.Sprint(want) {
		t.Fatalf("body = %s, want %v", rec.Body.String(), want)
	}
}
