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
)

func TestSwapSmokeSetupFlow(t *testing.T) {
	validator.SetupValidator() // normally done by httpserver.SetupFUN at startup
	server := handlers.NewServer(&services.Operations{})
	r := newTestRouter(t, server, middlewares{
		jwtAuth:    mwJWT,
		apiKeyAuth: mwJWT,
		anyAuth:    mwAnyAuth,
		clientOnly: mwClientOnly,
	})

	globals.MarkSetupComplete()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/auth/setup", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("GET /auth/setup (set up) -> %d %s", rec.Code, rec.Body.String())
	if rec.Code != http.StatusConflict {
		t.Fatalf("want 409 when set up, got %d", rec.Code)
	}

	// POST /auth/setup with an invalid body → validation middleware rejects
	req = httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/setup", strings.NewReader(`{"email":"nope","password":"weak"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("POST /auth/setup (invalid body) -> %d %s", rec.Code, rec.Body.String())
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for invalid body, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "VALIDATION_ERROR") {
		t.Fatalf("want VALIDATION_ERROR envelope, got %s", rec.Body.String())
	}

	// unknown route
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/nope", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("GET /nope -> %d", rec.Code)
}
