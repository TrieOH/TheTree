// Package httpserver is the HTTP-serving harness shared by every TrieOH backend:
// server lifecycle, pprof, middleware configuration, and the standard router skeleton
// (middleware stack, /health, OpenTelemetry wrapping). Telemetry is pushed
// over OTLP (see lib/telemetry); nothing is scraped from the process.
package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/pprof"
	"strings"
	"sync"
	"time"

	"lib/errx"
	"lib/telemetry"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// Config carries everything a backend must tell the harness about itself.
type Config struct {
	AppName            string
	Port               string
	ProfilePort        string
	CorsAllowedOrigins string
	CorsAllowedHeaders string
	// SkipLogPrefixes are extra request prefixes the request logger should
	// ignore, on top of the always-skipped /health.
	SkipLogPrefixes []string
	// RateLimit configures the per-client rate limiter. The zero value
	// keeps the historical defaults (400 rps, burst 20).
	RateLimit RateLimit
	// OpenAPISpec, when set, is served at GET /docs/openapi.yml — the
	// service's OpenAPI 3.1 specification (the same content as
	// api/<svc>/api-spec.yml), for docs tooling to fetch and render.
	OpenAPISpec []byte
	// Routes receives the chi router and registers the backend's feature
	// routes, auth middlewares, and any extra mounts (riverui, websockets).
	Routes func(r *chi.Mux)
}

// RateLimit is the per-client (real IP) token bucket. It limits per process:
// with several instances each has its own bucket, so an edge rate limit is
// the global limit and this one is a backstop.
type RateLimit struct {
	Disabled bool
	RPS      float64
	Burst    int
}

func (rl RateLimit) withDefaults() RateLimit {
	if rl.RPS <= 0 {
		rl.RPS = 400
	}
	if rl.Burst <= 0 {
		rl.Burst = 20
	}
	return rl
}

// cleanCSV splits a comma-separated env value, trims whitespace, and drops
// empty entries.
func cleanCSV(s string) []string {
	out := make([]string, 0)
	for part := range strings.SplitSeq(s, ",") {
		if v := strings.TrimSpace(part); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// health is the liveness probe: if it answers, the process is alive.
func health(service string) http.HandlerFunc {
	body, _ := json.Marshal(struct {
		Status  string `json:"status"`
		Service string `json:"service"`
	}{"ok", service})
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}

// spanRouteNameMiddleware renames the root span to "METHOD /route-template"
// once chi finishes routing, and records the matched route. otelhttp only
// renames a span when r.Pattern is set, which net/http's ServeMux fills in
// but chi never does — without this every trace is just named "GET"/"PATCH".
// Registered outermost so it runs even when a deeper middleware short-circuits.
func spanRouteNameMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		if span := trace.SpanFromContext(r.Context()); span.IsRecording() {
			if pattern := chi.RouteContext(r.Context()).RoutePattern(); pattern != "" {
				span.SetName(r.Method + " " + pattern)
				span.SetAttributes(semconv.HTTPRoute(pattern))
			} else {
				// No route matched (404): mirror the metrics middleware's
				// "not_found" label so untracked traffic is searchable
				// (resolveRoutePattern returns "not_found" for empty patterns).
				span.SetName(r.Method + " not_found")
			}
		}
	})
}

// NewRouter builds the standard router skeleton: chi with the standard
// middleware stack, the backend's routes, /health, wrapped in OpenTelemetry
// instrumentation that skips /health and OPTIONS.
func NewRouter(cfg Config) http.Handler {
	r := chi.NewRouter()

	r.Use(spanRouteNameMiddleware)
	for _, mw := range stack(cfg) {
		r.Use(mw)
	}

	if cfg.Routes != nil {
		cfg.Routes(r)
	}
	r.Get("/health", health(cfg.AppName))

	if len(cfg.OpenAPISpec) > 0 {
		spec := cfg.OpenAPISpec
		r.Get("/docs/openapi.yml", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/yaml")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(spec)
		})
	}

	return otelhttp.NewHandler(r, "http.server",
		// Name spans from the start (method + path) so even 404s/unmatched
		// routes are distinguishable; spanRouteNameMiddleware upgrades the
		// name to the route template once chi has matched it.
		otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/health" && r.URL.Path != "/docs/openapi.yml"
		}),
		// The river queue UI is a dev-only surface: don't trace it.
		otelhttp.WithFilter(func(r *http.Request) bool {
			return !strings.HasPrefix(r.URL.Path, "/riverui")
		}),
		// CORS preflights are noise: don't trace OPTIONS at all.
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.Method != http.MethodOptions
		}),
	)
}

// Start serves handler on cfg.Port, blocking until the process exits. It also
// starts the pprof server on cfg.ProfilePort when set.
//
// Superseded by Boot, which owns the whole process lifecycle.
func newServer(handler http.Handler, port string) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func servePprof(port, appName string) {
	pmux := http.NewServeMux()
	pmux.HandleFunc("/debug/pprof/", pprof.Index)
	pmux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pmux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pmux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	pmux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      pmux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Printf("%s pprof listening on :%s", appName, port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("%s pprof server error: %v", appName, err)
	}
}

var (
	collectorsOnce sync.Once
	collectorsVal  *httpCollectors
	errCollectors  error
)

func collectors() (*httpCollectors, error) {
	collectorsOnce.Do(func() {
		collectorsVal, errCollectors = newCollectors(prometheus.DefaultRegisterer)
	})
	return collectorsVal, errCollectors
}

// stack is the standard middleware stack shared by every backend.
func stack(cfg Config) []func(http.Handler) http.Handler {
	skipLog := append([]string{"/health", "/docs/openapi.yml"}, cfg.SkipLogPrefixes...)

	collectors, err := collectors()
	if err != nil {
		errx.Exit(err, "Failed to create collectors")
	}

	mwStack := []func(http.Handler) http.Handler{
		realIP,
		requestID,
		logs(telemetry.Log(), skipLog),
		metrics(collectors, []string{"/health"}),
		recoverer(telemetry.Log()),
		timeout(60 * time.Second),
		maxBodySize(1 << 20),
	}
	if rl := cfg.RateLimit.withDefaults(); !rl.Disabled {
		mwStack = append(mwStack, newRateLimiter(rl).middleware)
	}
	return append(mwStack, CORS(cfg.CorsAllowedOrigins, cfg.CorsAllowedHeaders))
}
