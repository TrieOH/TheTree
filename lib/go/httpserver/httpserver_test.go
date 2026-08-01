package httpserver

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func testRouter(t *testing.T, cfg Config) http.Handler {
	t.Helper()
	if cfg.AppName == "" {
		cfg.AppName = "test"
	}
	if cfg.Routes == nil {
		cfg.Routes = func(_ *chi.Mux) {}
	}
	return NewRouter(cfg)
}

func TestHealthReturnsAppName(t *testing.T) {
	h := testRouter(t, Config{AppName: "informd"})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("health status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "informd") {
		t.Fatalf("health body %q does not contain app name", rec.Body.String())
	}
}

func TestMetricsEndpoint(t *testing.T) {
	h := testRouter(t, Config{
		Routes: func(r *chi.Mux) {
			r.Get("/work", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
		},
	})

	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/work", nil))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/metrics", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("metrics status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "http_requests_total") {
		t.Fatalf("metrics body does not contain http_requests_total")
	}
}

func TestRoutesCallbackRegistered(t *testing.T) {
	h := testRouter(t, Config{
		Routes: func(r *chi.Mux) {
			r.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})
		},
	})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/ping", nil))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
}

func TestRequestIDHeaderSet(t *testing.T) {
	h := testRouter(t, Config{})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil))

	if rec.Header().Get("X-Request-ID") == "" {
		t.Fatalf("X-Request-ID not set on response")
	}
}

func TestCORSHeaders(t *testing.T) {
	h := testRouter(t, Config{CorsAllowedOrigins: "https://app.trieoh.com", CorsAllowedHeaders: "Content-Type"})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "https://app.trieoh.com")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatalf("Access-Control-Allow-Origin not set")
	}
}

// TestCORSPreflightAllowsConfiguredHeaders pins the contract the backends
// rely on: a header configured in CorsAllowedHeaders (e.g. Refresh-Token
// for logout/refresh) must come back in Access-Control-Allow-Headers, or
// the browser blocks the request.
func TestCORSPreflightAllowsConfiguredHeaders(t *testing.T) {
	h := testRouter(t, Config{
		CorsAllowedOrigins: "https://app.trieoh.com",
		CorsAllowedHeaders: "Content-Type,X-Request-ID,Authorization,Refresh-Token,X-API-Key",
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodOptions, "/auth/logout", nil)
	req.Header.Set("Origin", "https://app.trieoh.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "authorization, refresh-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	allowed := rec.Header().Get("Access-Control-Allow-Headers")
	for _, want := range []string{"Refresh-Token", "Authorization"} {
		found := false
		for _, h := range strings.Split(allowed, ",") {
			if strings.EqualFold(strings.TrimSpace(h), want) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Access-Control-Allow-Headers %q missing %q", allowed, want)
		}
	}
}

// TestCORSPreflightRejectsUnconfiguredHeaders pins the failure mode the
// backends hit before the CORS wiring landed: a requested header that is
// not in the allowed list must not be echoed back.
func TestCORSPreflightRejectsUnconfiguredHeaders(t *testing.T) {
	h := testRouter(t, Config{
		CorsAllowedOrigins: "https://app.trieoh.com",
		CorsAllowedHeaders: "Content-Type",
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodOptions, "/auth/logout", nil)
	req.Header.Set("Origin", "https://app.trieoh.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "refresh-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Headers"); strings.Contains(got, "Refresh-Token") {
		t.Fatalf("Refresh-Token must not be allowed when unconfigured, got %q", got)
	}
}

func TestRecoverReturnsError(t *testing.T) {
	h := testRouter(t, Config{
		Routes: func(r *chi.Mux) {
			r.Get("/boom", func(http.ResponseWriter, *http.Request) {
				panic("boom")
			})
		},
	})

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/boom", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
}

func TestBodySizeLimited(t *testing.T) {
	h := testRouter(t, Config{
		Routes: func(r *chi.Mux) {
			r.Post("/echo", func(w http.ResponseWriter, r *http.Request) {
				_, err := io.Copy(io.Discard, r.Body)
				if err != nil {
					w.WriteHeader(http.StatusRequestEntityTooLarge)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})
		},
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/echo", strings.NewReader(strings.Repeat("x", 2<<20)))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestNewServerTimeouts(t *testing.T) {
	srv := newServer(http.NewServeMux(), "8080")

	if srv.Addr != ":8080" {
		t.Fatalf("addr = %q, want :8080", srv.Addr)
	}
	if srv.ReadTimeout != 30*time.Second || srv.WriteTimeout != 60*time.Second || srv.IdleTimeout != 120*time.Second {
		t.Fatalf("timeouts = %v/%v/%v, want 30s/60s/120s", srv.ReadTimeout, srv.WriteTimeout, srv.IdleTimeout)
	}
}

func TestSetupFUNRuns(t *testing.T) {
	_ = t
	SetupFUN("test")
}
