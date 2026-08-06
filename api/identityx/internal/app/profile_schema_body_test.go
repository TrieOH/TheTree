package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"IdentityX/internal/handlers"
	"IdentityX/internal/services"
	"lib/globals"
	"lib/validator"

	"github.com/google/uuid"
)

// schemaTestRouter mounts the real strict middleware stack over an empty
// Operations set: bodies that decode go through to the handler and die at
// authn (no identity in context), while malformed bodies are rejected
// before the handler ever runs.
func schemaTestRouter(t *testing.T) http.Handler {
	t.Helper()
	validator.SetupValidator()
	globals.MarkSetupComplete()
	server := handlers.NewServer(&services.Operations{})
	return newTestRouter(t, server, middlewares{
		jwtAuth:    mwJWT,
		apiKeyAuth: mwJWT,
		anyAuth:    mwAnyAuth,
	})
}

func putProjectSchema(t *testing.T, r http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPut,
		"/projects/"+uuid.NewString()+"/profile-schema",
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

// TestUpsertProjectProfileSchemaAcceptsObjectSchema pins the regression
// where the spec forced the schema field to []byte: a JSON Schema object
// cannot unmarshal into []byte, so every valid body was rejected at decode
// with a generic "invalid request body" before reaching the handler. The
// body must now flow to the handler — here it is stopped at authn (401)
// because the test carries no identity, which proves the decode passed.
// The schema is the exact draft-2020-12 document a real client sent.
func TestUpsertProjectProfileSchemaAcceptsObjectSchema(t *testing.T) {
	r := schemaTestRouter(t)
	rec := putProjectSchema(t, r, `{"schema":{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"object","properties":{"full_name":{"type":"string","title":"Full name"},"display_name":{"type":"string","title":"Display name"},"picture_url":{"type":"string","format":"uri","title":"Profile picture"}},"required":["full_name"],"additionalProperties":false},"active":true}`)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("valid schema object must decode and reach the handler (authn rejects it): got %d %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "invalid request body") {
		t.Fatalf("valid schema object must not be reported as an invalid body, got %s", rec.Body.String())
	}
}

// TestUpsertProjectProfileSchemaMalformedJSONCarriesDetail pins that body
// decode failures surface the underlying cause in the envelope message
// instead of a bare "invalid request body".
func TestUpsertProjectProfileSchemaMalformedJSONCarriesDetail(t *testing.T) {
	r := schemaTestRouter(t)
	rec := putProjectSchema(t, r, `{"schema": ,"active":true}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for malformed JSON, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "invalid character") {
		t.Fatalf("want the decode cause in the message, got %s", rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), `"invalid request body"`) {
		t.Fatalf("want a detailed message, got a bare invalid request body: %s", rec.Body.String())
	}
}

// TestUpsertProjectProfileSchemaWrongFieldTypeCarriesDetail pins that a
// type mismatch in a sibling field (here: active must be a boolean) is
// reported with the offending field, not swallowed into a generic message.
func TestUpsertProjectProfileSchemaWrongFieldTypeCarriesDetail(t *testing.T) {
	r := schemaTestRouter(t)
	rec := putProjectSchema(t, r, `{"schema":{"type":"object"},"active":"yes"}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for wrong field type, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "cannot unmarshal") || !strings.Contains(rec.Body.String(), "active") {
		t.Fatalf("want the unmarshal cause naming the field, got %s", rec.Body.String())
	}
}
