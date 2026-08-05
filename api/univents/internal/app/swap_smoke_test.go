package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"univents/internal/handlers"
	"univents/internal/services"

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

	// JWT-protected route without a token -> 401 fun envelope via the dispatch
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/events/owned", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("GET /events/owned (no token) -> %d %s", rec.Code, rec.Body.String())
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
