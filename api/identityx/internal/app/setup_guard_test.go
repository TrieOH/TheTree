package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"IdentityX/internal/setup"
	"lib/validator"
)

// TestSetupGuardBlocksUntilSetup pins the guard: before setup completes,
// every operation but the /auth/setup pair answers 503.
func TestSetupGuardBlocksUntilSetup(t *testing.T) {
	validator.SetupValidator()
	st := setup.NewState(func(context.Context) (bool, error) { return false, nil })
	r := newTestRouter(t, newServer(operations{}, st), testPrimitives(), st)

	// Required params bind before the chains run, so the request carries
	// both logout headers; the token lacks "Bearer " so the handler answers
	// 401 after setup without reaching the nil operations.
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
