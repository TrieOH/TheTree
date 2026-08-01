package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"IdentityX/internal/handlers"
	"IdentityX/internal/services"
	"lib/globals"
)

func TestSwapSmokeSetupFlow(t *testing.T) {
	server := handlers.NewServer(&services.Operations{})
	r := newTestRouter(server, middlewares{
		jwtAuth:    mwJWT,
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

	// POST /auth/setup with bad body shape → validation middleware
	req = httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/auth/setup", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("POST /auth/setup (no body) -> %d", rec.Code)

	// unknown route
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/nope", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("GET /nope -> %d", rec.Code)
}
