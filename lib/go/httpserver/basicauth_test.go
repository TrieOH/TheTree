package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuthGatesTheSubtree(t *testing.T) {
	t.Setenv(basicAuthUserEnv, "ops")
	t.Setenv(basicAuthPassEnv, "hunter2")

	h := BasicAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("no credentials is 401 with challenge", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rec.Code)
		}
		if rec.Header().Get("WWW-Authenticate") == "" {
			t.Fatal("missing WWW-Authenticate challenge")
		}
	})

	t.Run("wrong credentials is 401", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		req.SetBasicAuth("ops", "wrong")
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rec.Code)
		}
	})

	t.Run("correct credentials passes", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		req.SetBasicAuth("ops", "hunter2")
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})
}

func TestBasicAuthFailsClosedWhenUnconfigured(t *testing.T) {
	t.Setenv(basicAuthUserEnv, "")
	t.Setenv(basicAuthPassEnv, "")

	h := BasicAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Even a correct-looking credential pair must not pass when the env is
	// unset: the endpoint stays closed rather than silently going public.
	rec := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	req.SetBasicAuth("ops", "hunter2")
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
}
