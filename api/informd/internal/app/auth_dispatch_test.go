package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type adKey string

const adSentinelKey adKey = "auth-dispatch-sentinel"

var adSentinel = &struct{}{}

// Regression: oapi-codegen passes the operationID in generated (PascalCase)
// form (e.g. "CreateNamespace"), while the chains map is keyed camelCase
// ("createNamespace"). If the casing is not normalized, the lookup misses and
// no auth middleware runs — every protected operation becomes public.
// Additionally, auth middlewares replace the request (and its context) via
// r.WithContext(...); the handler must observe that modified context, not
// the pre-middleware one.
func TestAuthDispatchRunsChainForProtectedOperation(t *testing.T) {
	var ran bool
	jwt := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ran = true
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), adSentinelKey, adSentinel)))
		})
	}
	dispatch := authDispatch(middlewares{jwt: jwt})

	handler := dispatch(func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		if got := ctx.Value(adSentinelKey); got != adSentinel {
			return nil, fmt.Errorf("handler saw stale context: sentinel %v, want %v", got, adSentinel)
		}
		return "ok", nil
	}, "CreateNamespace") // exactly as oapi-codegen passes it

	resp, err := handler(context.Background(), httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/namespaces", nil), nil)
	if err != nil {
		t.Fatalf("handler must see the middleware-modified request context: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("want %q, got %v", "ok", resp)
	}
	if !ran {
		t.Fatal("auth middleware never ran: operationID casing mismatch")
	}
}

func TestAuthDispatchSkipsPublicOperation(t *testing.T) {
	var ran bool
	dispatch := authDispatch(middlewares{jwt: func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ran = true
			next.ServeHTTP(w, r)
		})
	}})
	handler := dispatch(func(_ context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		return "ok", nil
	}, "GetAnswerableForm")
	_, err := handler(context.Background(), httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/forms/{formID}/answerable", nil), nil)
	if err != nil {
		t.Fatalf("public operation must pass through: %v", err)
	}
	if ran {
		t.Fatal("public operation must not run the auth chain")
	}
}

func TestAuthDispatchRejectionShortCircuits(t *testing.T) {
	dispatch := authDispatch(middlewares{jwt: func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		})
	}})
	var called bool
	handler := dispatch(func(_ context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		called = true
		return "ok", nil
	}, "CreateNamespace")
	rec := httptest.NewRecorder()
	resp, err := handler(context.Background(), rec, httptest.NewRequest(http.MethodPost, "/namespaces", nil), nil)
	if err != nil {
		t.Fatalf("rejected request must return nil error: %v", err)
	}
	if resp != nil {
		t.Fatalf("rejected request must return nil response, got %v", resp)
	}
	if called {
		t.Fatal("handler must not run after auth rejection")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rec.Code)
	}
}
