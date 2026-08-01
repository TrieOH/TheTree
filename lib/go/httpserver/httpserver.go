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

// NewRouter builds the standard router skeleton: chi with the standard
// middleware stack, /metrics, the backend's routes, /health, wrapped in
// OpenTelemetry instrumentation that skips /metrics and /health.
func NewRouter(cfg Config) http.Handler {
	r := chi.NewRouter()

	for _, mw := range stack(cfg) {
		r.Use(mw)
	}

	r.Handle("/metrics", promhttp.Handler())

	if cfg.Routes != nil {
		cfg.Routes(r)
	}
	r.Get("/health", fh.Health(cfg.AppName).Handle)

	return otelhttp.NewHandler(r, "http.server",
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/health"
		}),
		otelhttp.WithFilter(func(r *http.Request) bool {
			return r.URL.Path != "/metrics"
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
	skipLog := append([]string{"/metrics", "/health"}, cfg.SkipLogPrefixes...)

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
