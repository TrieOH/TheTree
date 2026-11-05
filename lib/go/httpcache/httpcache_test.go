package httpcache

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func serve(h http.Handler) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	Middleware(h).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	return rec
}

func TestDefaultIsPrivate(t *testing.T) {
	rec := serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Tag(r.Context(), "event:1")
		_, _ = w.Write([]byte("ok"))
	}))
	if got := rec.Header().Get("Cache-Control"); got != Private {
		t.Fatalf("Cache-Control = %q, want %q", got, Private)
	}
	if got := rec.Header().Get(TagHeader); got != "" {
		t.Fatalf("a private response carries tags %q", got)
	}
}

func TestSharedOK(t *testing.T) {
	rec := serve(Shared(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Tag(r.Context(), "events", "event:1", "events")
		w.WriteHeader(http.StatusOK)
	})))
	h := rec.Header()
	if h.Get("Cache-Control") != SharedClient || h.Get("CDN-Cache-Control") != SharedEdge {
		t.Fatalf("headers = %v", h)
	}
	if got := h.Get(TagHeader); got != "events,event:1" {
		t.Fatalf("%s = %q", TagHeader, got)
	}
}

func TestSharedErrorIsPrivate(t *testing.T) {
	rec := serve(Shared(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Tag(r.Context(), "event:1")
		w.WriteHeader(http.StatusNotFound)
	})))
	if got := rec.Header().Get("Cache-Control"); got != Private {
		t.Fatalf("Cache-Control = %q, want %q", got, Private)
	}
	if rec.Header().Get("CDN-Cache-Control") != "" || rec.Header().Get(TagHeader) != "" {
		t.Fatalf("a 404 is marked shared: %v", rec.Header())
	}
}

func TestPurgeOnlyOnSuccess(t *testing.T) {
	for code, want := range map[int]string{http.StatusOK: "a,b", http.StatusNoContent: "a,b", http.StatusConflict: ""} {
		rec := serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Purge(r.Context(), "a")
			Purge(r.Context(), "b", "a")
			w.WriteHeader(code)
		}))
		if got := rec.Header().Get(PurgeHeader); got != want {
			t.Fatalf("status %d: %s = %q, want %q", code, PurgeHeader, got, want)
		}
	}
}

func TestHandlerPolicyWins(t *testing.T) {
	rec := serve(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		Deploy(w)
		_, _ = w.Write([]byte("spec"))
	}))
	if rec.Header().Get("Cache-Control") != SharedClient || rec.Header().Get("CDN-Cache-Control") != DeployEdge {
		t.Fatalf("headers = %v", rec.Header())
	}
}

func TestNoMiddlewareIsNoop(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	Tag(r.Context(), "x")
	Purge(r.Context(), "x")
}
