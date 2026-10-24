package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	spec "IdentityX"
	"IdentityX/internal/config"
	"IdentityX/internal/handlers"
	"IdentityX/internal/services"
	"IdentityX/internal/setup"
	"lib/backend"
	libconfig "lib/config"
	"lib/httpserver"
	"lib/validator"

	"github.com/go-chi/chi/v5"
)

// TestCORSWiringThroughTheKit pins that the CORS values read from env
// (config) reach the router frame: a preflight requesting the Refresh-Token
// header must be answered with it allowed, or the browser blocks
// logout/refresh.
func TestCORSWiringThroughTheKit(t *testing.T) {
	validator.SetupValidator()
	cfg := config.Config{
		CORS: libconfig.CORS{
			AllowedOrigins: "http://localhost:3000",
			AllowedHeaders: "Content-Type,X-Request-ID,Authorization,Refresh-Token,X-API-Key",
		},
	}
	srv := backend.ServerFrom(cfg.Server, cfg.CORS, "")
	server := handlers.NewServer(&services.Operations{}, completedSetup())
	r := httpserver.NewRouter(httpserver.Config{
		CorsAllowedOrigins: srv.CorsAllowedOrigins,
		CorsAllowedHeaders: srv.CorsAllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		RateLimit:          httpserver.RateLimit{Disabled: true},
		Routes: func(r *chi.Mux) {
			mountStrict(r, server, newTestRouterChains(t, testPrimitives(), completedSetup()))
		},
	})

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
		t.Fatalf("Access-Control-Allow-Headers %q missing Refresh-Token — CORS config not wired to the router frame", allowed)
	}
}

// TestSetupGuardBlocksUntilSetup pins the guard: before setup completes,
// every operation but the /auth/setup pair answers 503.
func TestSetupGuardBlocksUntilSetup(t *testing.T) {
	validator.SetupValidator()
	st := setup.NewState(func(context.Context) (bool, error) { return false, nil })
	r := newTestRouter(t, handlers.NewServer(&services.Operations{}, st), testPrimitives(), st)

	// Required params bind before the chains run, so the request carries
	// both logout headers; the token lacks "Bearer " so the handler answers
	// 401 after setup without reaching the nil services.
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Authorization", "token")
	req.Header.Set("Refresh-Token", "token")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("want 503 before setup, got %d %s", rec.Code, rec.Body.String())
	}

	st.MarkComplete()
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code == http.StatusServiceUnavailable {
		t.Fatalf("guard still blocking after setup: %s", rec.Body.String())
	}
}

func TestSwapSmokeSetupFlow(t *testing.T) {
	validator.SetupValidator() // normally done by httpserver.SetupFUN at startup
	st := completedSetup()
	r := newTestRouter(t, handlers.NewServer(&services.Operations{}, st), testPrimitives(), st)

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
