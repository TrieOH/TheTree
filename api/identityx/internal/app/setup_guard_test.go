package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
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

// runChain serves one request through chain and reports whether it reached
// the handler.
func runChain(chain []func(http.Handler) http.Handler) (reached bool) {
	var next http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		reached = true
		w.WriteHeader(http.StatusNoContent)
	})
	for _, mw := range slices.Backward(chain) {
		next = mw(next)
	}
	next.ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))
	return reached
}

func TestGuardSetupBlocksEverythingButTheSetupRoutes(t *testing.T) {
	chains := map[string][]func(http.Handler) http.Handler{
		"GetSetup":    nil,
		"PostSetup":   nil,
		"ListActors":  nil,
		"GetProjects": {mwJWT},
	}
	err := guardSetup(chains, setup.NewState(nil)) // setup not done
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"GetSetup", "PostSetup"} {
		if reached := runChain(chains[id]); !reached {
			t.Errorf("%s must skip the guard", id)
		}
	}
	for _, id := range []string{"ListActors", "GetProjects"} {
		if reached := runChain(chains[id]); reached {
			t.Errorf("%s must be refused until setup is done", id)
		}
	}

	done := map[string][]func(http.Handler) http.Handler{"GetSetup": nil, "PostSetup": nil, "ListActors": nil}
	err = guardSetup(done, completedSetup())
	if err != nil {
		t.Fatal(err)
	}
	if reached := runChain(done["ListActors"]); !reached {
		t.Error("ListActors must pass once setup is done")
	}
}

func TestGuardSetupFailsWithoutTheSetupRoutes(t *testing.T) {
	err := guardSetup(map[string][]func(http.Handler) http.Handler{"GetSetup": nil}, completedSetup())
	if err == nil || !strings.Contains(err.Error(), "postSetup") {
		t.Fatalf("want an error naming postSetup, got %v", err)
	}
}
