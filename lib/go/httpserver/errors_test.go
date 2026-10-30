package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"lib/database"
	"lib/fielderr"

	fp "github.com/MintzyG/problematic"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

const testAPIURL = "https://api.example.test"

var (
	errAddSelf      = errors.New("users can't add themselves to projects")
	errTermsMissing = errors.New("you must accept the terms of service to register")
)

// testMapper stands in for a Backend's generated problems.Mapper.
func testMapper(err error) *fp.Problem {
	var notFound *database.NotFoundError
	switch {
	case errors.Is(err, errAddSelf):
		return fp.New(http.StatusBadRequest).WithType("add-self").WithTitle("Add Self").WithDetail(errAddSelf.Error())
	case errors.Is(err, errTermsMissing):
		return fp.New(http.StatusBadRequest).WithType("validation").WithTitle("Validation Failed").WithDetail(errTermsMissing.Error())
	case errors.As(err, &notFound):
		return fp.New(http.StatusNotFound).WithType("not-found").WithTitle("Not Found").WithDetail(notFound.Error())
	}
	return nil
}

// errBadKey is a test authenticator's rejection.
var errBadKey = errors.New("invalid api key")

// harnessMapper stands in for the generated shared cases: the harness's
// own rejections, as every Backend answers them.
func harnessMapper(err error) *fp.Problem {
	switch {
	case errors.Is(err, ErrMissingBearer), errors.Is(err, ErrMissingAPIKey), errors.Is(err, ErrNoCredentials), errors.Is(err, errBadKey):
		return fp.New(http.StatusUnauthorized).WithType("unauthenticated").WithDetail(err.Error())
	case errors.Is(err, ErrRateLimited):
		return fp.New(http.StatusTooManyRequests).WithType("rate-limited").WithDetail(ErrRateLimited.Error())
	case errors.Is(err, ErrInvalidParameter), errors.Is(err, ErrInvalidHeader), errors.As(err, new(*InvalidBodyError)):
		return fp.New(http.StatusBadRequest).WithType("validation").WithDetail(err.Error())
	}
	return nil
}

// TestMain installs harnessMapper, as a Backend's boot installs its mapper.
func TestMain(m *testing.M) {
	UseProblems(harnessMapper, testAPIURL)
	os.Exit(m.Run())
}

// useTestProblems installs m for one test and puts harnessMapper back
// after it.
func useTestProblems(t *testing.T, m fp.ErrorMapper) {
	t.Helper()
	UseProblems(m, testAPIURL+"/")
	t.Cleanup(func() { UseProblems(harnessMapper, testAPIURL) })
}

// writeProblem runs WriteError and decodes the problem it wrote.
func writeProblem(t *testing.T, path string, err error) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, path, nil)
	WriteError(rec, req, err)
	if ct := rec.Header().Get("Content-Type"); ct != fp.ContentType {
		t.Fatalf("Content-Type = %q, want %q", ct, fp.ContentType)
	}
	var body map[string]any
	decodeErr := json.Unmarshal(rec.Body.Bytes(), &body)
	if decodeErr != nil {
		t.Fatalf("decode: %v", decodeErr)
	}
	return rec, body
}

// TestWriteErrorIsAProblem — a mapped error goes out as its problem, the
// type resolved against API_URL and the request path as instance.
func TestWriteErrorIsAProblem(t *testing.T) {
	useTestProblems(t, testMapper)

	rec, body := writeProblem(t, "/projects/x/members", fmt.Errorf("add member: %w", errAddSelf))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
	want := map[string]any{
		"type":     testAPIURL + "/problems/add-self",
		"status":   float64(400),
		"title":    "Add Self",
		"detail":   errAddSelf.Error(),
		"instance": "/projects/x/members",
	}
	if fmt.Sprint(body) != fmt.Sprint(want) {
		t.Fatalf("body = %s, want %v", rec.Body.String(), want)
	}
}

// TestWriteErrorChecksXProblems — a problem the answering operation's
// x-problems doesn't list panics under test; a listed one, or an answer
// outside an operation, goes out as is.
func TestWriteErrorChecksXProblems(t *testing.T) {
	useTestProblems(t, testMapper)
	err := UseOperationProblems([]byte(`
paths:
  /members:
    post:
      operationId: addMember
      x-problems: [add-self]
`), false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = UseOperationProblems(nil, false) })

	// serve answers like the generated strict server: the handler's error
	// is written with the request the server was handed.
	serve := func(answer error, operationID string) func() {
		return func() {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				f := operationMiddleware()(func(context.Context, http.ResponseWriter, *http.Request, any) (any, error) {
					return nil, answer
				}, operationID)
				_, err := f(r.Context(), w, r, nil)
				WriteError(w, r, err)
			})
			req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/members", nil)
			OperationSlot(handler).ServeHTTP(httptest.NewRecorder(), req)
		}
	}
	mustNotPanic := func(name string, f func()) {
		t.Helper()
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("%s panicked: %v", name, r)
			}
		}()
		f()
	}

	mustNotPanic("listed", serve(errAddSelf, "AddMember"))
	mustNotPanic("undeclared operation", serve(errTermsMissing, "Other"))
	mustNotPanic("outside an operation", func() { writeProblem(t, "/members", errTermsMissing) })
	defer func() {
		if r := recover(); r == nil || !strings.Contains(fmt.Sprint(r), "AddMember answered problem validation") {
			t.Errorf("unlisted: recover() = %v, want the unlisted problem named", r)
		}
	}()
	serve(errTermsMissing, "AddMember")()
}

// TestWriteErrorUnmappedIsGeneric — an error the mapper doesn't know, or
// any error before a mapper is installed, is a 500 that never shows its
// cause.
func TestWriteErrorUnmappedIsGeneric(t *testing.T) {
	for name, m := range map[string]fp.ErrorMapper{"unmapped": testMapper, "no mapper": nil} {
		t.Run(name, func(t *testing.T) {
			useTestProblems(t, m)
			rec, body := writeProblem(t, "/projects", errors.New("pq: connection reset"))
			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("want 500, got %d", rec.Code)
			}
			if body["detail"] != "Internal Server Error" || body["type"] != testAPIURL+"/problems/internal-server-error" {
				t.Fatalf("body = %s", rec.Body.String())
			}
		})
	}
}

// TestWriteErrorListsFields — a fielderr.Lister in the chain names its
// fields in the problem's errors member.
func TestWriteErrorListsFields(t *testing.T) {
	useTestProblems(t, testMapper)

	_, body := writeProblem(t, "/auth/register", fielderr.On(errTermsMissing, "accept_terms", "required"))

	want := []any{map[string]any{"field": "accept_terms", "message": "required"}}
	if fmt.Sprint(body["errors"]) != fmt.Sprint(want) {
		t.Fatalf("errors = %v, want %v", body["errors"], want)
	}
}

// TestTraceMiddlewareLeavesClientErrorsUnset — a 4xx is recorded on the
// span but does not mark it failed.
func TestTraceMiddlewareLeavesClientErrorsUnset(t *testing.T) {
	useTestProblems(t, testMapper)
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

// TestProblemDocs — GET /problems/{slug} documents a declared problem with
// an example body; an unknown slug is the not-found problem.
func TestProblemDocs(t *testing.T) {
	useTestProblems(t, testMapper)
	r := chi.NewRouter()
	r.Get("/problems/{slug}", serveProblemDocs([]ProblemDoc{{
		Slug: "payment-rejected", Status: 400, Title: "Payment Rejected",
		Description: "The provider refused the payment.", Extensions: json.RawMessage(`{"message":""}`),
	}}))

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/problems/payment-rejected", nil))
	var doc struct {
		Type    string         `json:"type"`
		Status  int            `json:"status"`
		Example map[string]any `json:"example"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &doc)
	if err != nil || rec.Code != http.StatusOK {
		t.Fatalf("got %d %s (%v)", rec.Code, rec.Body.String(), err)
	}
	if doc.Type != testAPIURL+"/problems/payment-rejected" || doc.Status != 400 {
		t.Fatalf("doc = %+v", doc)
	}
	if _, ok := doc.Example["message"]; !ok || doc.Example["type"] != doc.Type {
		t.Fatalf("example = %v", doc.Example)
	}

	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/problems/nope", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown slug: got %d", rec.Code)
	}
}
