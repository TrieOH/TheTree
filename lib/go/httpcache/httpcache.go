// Package httpcache decides how shared caches may store a response, in
// standard headers any CDN can read. Every response is private, no-store
// unless its operation opts in: an operation the spec marks x-cache: shared
// (ADR-0020) may be kept at the edge until a write purges one of its tags.
//
// Handlers and services name what a response holds with Tag and what a
// write changed with Purge; both read the request context, so they are
// no-ops where no HTTP request is in flight (jobs, tests).
package httpcache

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"strings"
)

// Response headers this package writes. TagHeader lists the tags a shared
// response carries; PurgeHeader lists the tags a successful write made
// stale. The edge consumes both and strips them before the client sees
// the response.
const (
	TagHeader   = "Cache-Tag"
	PurgeHeader = "Cache-Purge-Tag"
)

// Policies, as Cache-Control and CDN-Cache-Control values.
const (
	// Private is every response's default.
	Private = "private, no-store"
	// SharedClient makes browsers revalidate every time: only the edge,
	// which can be purged, keeps a shared response.
	SharedClient = "public, max-age=0, must-revalidate"
	// SharedEdge is how long the edge may keep a shared response. A purge
	// normally ends it first; the TTL only heals a purge that was lost.
	SharedEdge = "max-age=604800"
	// DeployEdge is how long the edge may keep a response that only
	// changes with a deploy. The edge's cache is per Worker version, so a
	// deploy starts it empty.
	DeployEdge = "max-age=31536000"
)

type state struct {
	shared bool
	tags   []string
	purge  []string
}

type stateKey struct{}

func from(ctx context.Context) *state {
	s, _ := ctx.Value(stateKey{}).(*state)
	return s
}

// Tag names what the response holds, so a write that changes any of it
// can purge the response. It only takes effect on a shared response.
func Tag(ctx context.Context, tags ...string) {
	if s := from(ctx); s != nil {
		s.tags = append(s.tags, tags...)
	}
}

// Purge names what a write changed. The edge purges these tags once the
// write succeeds (any 2xx); a failed write purges nothing.
func Purge(ctx context.Context, tags ...string) {
	if s := from(ctx); s != nil {
		s.purge = append(s.purge, tags...)
	}
}

// Shared marks the request's response as storable at the edge. It is the
// per-operation middleware the spec's x-cache: shared resolves to; only a
// 200 is stored.
func Shared(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s := from(r.Context()); s != nil {
			s.shared = true
		}
		next.ServeHTTP(w, r)
	})
}

// Deploy marks a response that only changes with a deploy (the spec, the
// problem docs). Call it before writing the body.
func Deploy(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", SharedClient)
	w.Header().Set("CDN-Cache-Control", DeployEdge)
}

// Middleware gives every request a policy and writes it when the response
// starts. A Cache-Control the handler set itself is left alone.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := &state{}
		ww := &writer{ResponseWriter: w, state: s}
		next.ServeHTTP(ww, r.WithContext(context.WithValue(r.Context(), stateKey{}, s)))
	})
}

type writer struct {
	http.ResponseWriter

	state   *state
	written bool
}

func (w *writer) apply(code int) {
	if w.written {
		return
	}
	w.written = true
	h := w.Header()
	if h.Get("Cache-Control") == "" {
		if w.state.shared && code == http.StatusOK {
			h.Set("Cache-Control", SharedClient)
			h.Set("CDN-Cache-Control", SharedEdge)
			if len(w.state.tags) > 0 {
				h.Set(TagHeader, join(w.state.tags))
			}
		} else {
			h.Set("Cache-Control", Private)
		}
	}
	if code >= 200 && code < 300 && len(w.state.purge) > 0 {
		h.Set(PurgeHeader, join(w.state.purge))
	}
}

func (w *writer) WriteHeader(code int) {
	w.apply(code)
	w.ResponseWriter.WriteHeader(code)
}

func (w *writer) Write(b []byte) (int, error) {
	w.apply(http.StatusOK)
	return w.ResponseWriter.Write(b)
}

func (w *writer) Flush() {
	w.apply(http.StatusOK)
	_ = http.NewResponseController(w.ResponseWriter).Flush()
}

// Hijack hands the connection over as is (the SSE streams write their own
// headers on it).
func (w *writer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w.written = true
	return http.NewResponseController(w.ResponseWriter).Hijack()
}

func (w *writer) Unwrap() http.ResponseWriter { return w.ResponseWriter }

// join dedupes tags, keeping their order.
func join(tags []string) string {
	seen := make(map[string]bool, len(tags))
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		if t != "" && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return strings.Join(out, ",")
}
