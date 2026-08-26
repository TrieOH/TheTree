// Package httpserver is the HTTP-serving harness shared by every TrieOH backend:
// server lifecycle, pprof, fun configuration, and the standard router skeleton
// (middleware stack, /metrics, /health, OpenTelemetry wrapping).
package httpserver

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"strings"
	"sync"
	"time"

	"lib/errx"
	"lib/telemetry"
	"lib/validator"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	fh "github.com/MintzyG/fun/handlers"
	mws "github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Config carries everything a backend must tell the harness about itself.
type Config struct {
	AppName            string
	Port               string
	ProfilePort        string
	CorsAllowedOrigins string
	CorsAllowedHeaders string
	// SkipLogPrefixes are extra request prefixes the request logger should
	// ignore, on top of the always-skipped /metrics and /health.
	SkipLogPrefixes []string
	// OpenAPISpec, when set, is served at GET /docs/openapi.yml — the
	// service's OpenAPI 3.1 specification (the same content as
	// api/<svc>/api-spec.yml), for docs tooling to fetch and render.
	OpenAPISpec []byte
	// Routes receives the chi router and registers the backend's feature
	// routes, auth middlewares, and any extra mounts (riverui, websockets).
	Routes func(r *chi.Mux)
}

// SetupFUN configures the shared fun runtime for the process. Call once at
// startup, before any request is handled.
func SetupFUN(appName string) {
	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        appName,
	})

	v := validator.SetupValidator()
	bind.SetValidator(v)
	fun.SetPathParamFunc(chi.URLParam)

	err := fun.AddInterceptor(&simpleLogger{})
	if err != nil {
		errx.Exit(err, "failed to add interceptor")
	}
}

type simpleLogger struct{}

func (l *simpleLogger) Intercept(_ context.Context, rs *fun.Response, statusCode int) {
	if statusCode == 500 {
		telemetry.Log().Info("InternalServerError Response", zap.Any("response", rs))
	}
}

func (l *simpleLogger) InterceptSimple(rs *fun.Response, statusCode int) {
	if statusCode == 500 {
		telemetry.Log().Info("InternalServerError Response", zap.Any("response", rs))
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
// middleware stack, /metrics, the backend's routes, /health, wrapped in
// OpenTelemetry instrumentation that skips /metrics, /health and OPTIONS.
func NewRouter(cfg Config) http.Handler {
	r := chi.NewRouter()

	r.Use(spanRouteNameMiddleware)
	for _, mw := range stack(cfg) {
		r.Use(mw)
	}

	r.Handle("/metrics", promhttp.Handler())

	if cfg.Routes != nil {
		cfg.Routes(r)
	}
	r.Get("/health", fh.Health(cfg.AppName).Handle)

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
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/metrics"
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
func Start(handler http.Handler, cfg Config) {
	if cfg.ProfilePort != "" {
		go servePprof(cfg.ProfilePort, cfg.AppName)
	}

	log.Printf("%s listening on :%s", cfg.AppName, cfg.Port)
	log.Fatal(newServer(handler, cfg.Port).ListenAndServe())
}

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
	collectorsVal  *mws.Collectors
	errCollectors  error
)

func collectors() (*mws.Collectors, error) {
	collectorsOnce.Do(func() {
		collectorsVal, errCollectors = mws.NewCollectors(prometheus.DefaultRegisterer)
	})
	return collectorsVal, errCollectors
}

// stack is the standard middleware stack shared by every backend.
func stack(cfg Config) []func(http.Handler) http.Handler {
	skipLog := append([]string{"/metrics", "/health", "/docs/openapi.yml"}, cfg.SkipLogPrefixes...)

	collectors, err := collectors()
	if err != nil {
		errx.Exit(err, "Failed to create collectors")
	}

	return []func(http.Handler) http.Handler{
		mws.RealIP(),
		mws.RequestID(mws.RequestIDConfig{Header: "X-Request-ID"}),
		mws.Logs(mws.Config{Logger: telemetry.Log(), SkipPrefixes: skipLog, RequestIDHeader: "X-Request-ID"}),
		mws.Metrics(collectors, mws.MetricsConfig{SkipPrefixes: []string{"/metrics", "/health"}}),
		mws.Recover(telemetry.Log()),
		mws.Timeout(60 * time.Second),
		mws.MaxBodySize(1 << 20),
		mws.RateLimit(mws.RateLimitConfig{
			RPS:   400,
			Burst: 20,
			KeyExtractor: func(r *http.Request) string {
				return r.RemoteAddr
			},
		}),
		mws.CORS(mws.CORSConfig{
			AllowedOrigins:   xslices.Clean(strings.Split(cfg.CorsAllowedOrigins, ",")),
			AllowedHeaders:   xslices.Clean(strings.Split(cfg.CorsAllowedHeaders, ",")),
			AllowCredentials: true,
		}),
	}
}
