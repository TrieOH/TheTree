package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"IdentityX/internal/config"
	"IdentityX/internal/handlers"
	"IdentityX/internal/services"
	"lib/globals"
	"lib/validator"
)

// TestCORSWiringThroughCreateRouter pins that the CORS values read from
// env (config) actually reach the harness: a preflight requesting the
// Refresh-Token header must be answered with it allowed, or the browser
// blocks logout/refresh.
func TestCORSWiringThroughCreateRouter(t *testing.T) {
	validator.SetupValidator()
	server := handlers.NewServer(&services.Operations{})
	app := &IdentityX{cfg: config.Config{
		CorsAllowedOrigins: "http://localhost:3000",
		CorsAllowedHeaders: "Content-Type,X-Request-ID,Authorization,Refresh-Token,X-API-Key",
	}}
	r := app.CreateRouter(middlewares{
		jwtAuth:    mwJWT,
		apiKeyAuth: mwJWT,
		anyAuth:    mwAnyAuth,
		scopes:     testScopeCheckers(),
	}, server, nil)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodOptions, "/auth/logout", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "authorization, refresh-token")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	allowed := rec.Header().Get("Access-Control-Allow-Headers")
	found := false
	for h := range strings.SplitSeq(allowed, ",") {
		if strings.EqualFold(strings.TrimSpace(h), "Refresh-Token") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Access-Control-Allow-Headers %q missing Refresh-Token — CORS config not wired to the harness", allowed)
	}
}

func TestSwapSmokeSetupFlow(t *testing.T) {
	validator.SetupValidator() // normally done by httpserver.SetupFUN at startup
	server := handlers.NewServer(&services.Operations{})
	r := newTestRouter(t, server, middlewares{
		jwtAuth:    mwJWT,
		apiKeyAuth: mwJWT,
		anyAuth:    mwAnyAuth,
		scopes:     testScopeCheckers(),
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
