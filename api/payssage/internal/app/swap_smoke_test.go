package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"payssage/internal/handlers"
	"payssage/internal/services"

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
	r := newTestRouter(server, middlewares{jwtAuth: rejectJWT})

	// JWT-protected route without a token -> 401 fun envelope via the dispatch
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/wallets", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("GET /wallets (no token) -> %d %s", rec.Code, rec.Body.String())
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rec.Code)
	}

	// public webhook receive still reaches the handler (rejectJWT not applied)
	req = httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/webhooks/mercado_pago", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("POST /webhooks/mercado_pago -> %d", rec.Code)

	// unknown route -> 404
	req = httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/nope", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	t.Logf("GET /nope -> %d", rec.Code)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}
}
