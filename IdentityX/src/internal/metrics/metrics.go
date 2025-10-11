package metrics

import (
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handler",
			Buckets: prometheus.ExponentialBuckets(0.00025, 2, 16),
			// starts at 250µs (0.00025s), multiplies by 2, 16 buckets:
			// 250µs, 500µs, 1000µs, 2.0ms, 4.0ms, ... up to ~8s
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(HttpRequestsTotal, HttpRequestDuration)
}

func MetricsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" || strings.HasPrefix(r.URL.Path, "/swagger") {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()

		ww := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()

		route := normalizePath(r)

		HttpRequestsTotal.WithLabelValues(route, r.Method, http.StatusText(ww.status)).Inc()
		HttpRequestDuration.WithLabelValues(route).Observe(duration)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
