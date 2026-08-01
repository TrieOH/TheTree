package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"Informd/internal/handlers"
	"Informd/internal/services"

	"github.com/MintzyG/fun"
)

// rejectJWT simulates a JWT middleware that always rejects.
func rejectJWT(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fun.Unauthorized("missing or invalid token").Send(w)
	})
}

func TestSwapSmokeAuthDispatch(t *testing.T) {
	server := handlers.NewServer(&services.Operations{})
	r := newTestRouter(t, server, middlewares{jwt: rejectJWT})

	// JWT-only namespace route without a token -> 401 fun envelope
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/namespaces", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("GET /namespaces (no token) -> %d %s", rec.Code, rec.Body.String())
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rec.Code)
	}

	// unknown route -> 404
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/nope", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("GET /nope -> %d", rec.Code)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}
}
