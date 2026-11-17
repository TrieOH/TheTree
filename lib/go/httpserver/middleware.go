package httpserver

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"runtime/debug"
	"strings"
	"sync"
	"time"
	"uuid"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/time/rate"
)

const requestIDHeader = "X-Request-ID"

type requestIDKey struct{}

// RequestIDFromCtx is the request's ID as the requestID middleware set it.
func RequestIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey{}).(string)
	return v
}

// realIP sets r.RemoteAddr to the client's bare IP: CF-Connecting-IP when
// present, else the peer's IP with its port stripped, so the rate limiter
// keys per client and not per connection. Backends run as Cloudflare
// Containers reached only through their Worker, and Cloudflare's edge
// overwrites CF-Connecting-IP on every request, so a client cannot pick
// it. X-Forwarded-For is never read: Cloudflare doesn't send it and any
// value would come from the client.
func realIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ip, ok := parseIP(r.Header.Get("CF-Connecting-IP")); ok {
			r.RemoteAddr = ip.String()
		} else if ip, ok := peerAddr(r.RemoteAddr); ok {
			r.RemoteAddr = ip.String()
		}
		next.ServeHTTP(w, r)
	})
}

func peerAddr(remoteAddr string) (netip.Addr, bool) {
	ap, err := netip.ParseAddrPort(remoteAddr)
	if err == nil {
		return ap.Addr().Unmap(), true
	}
	return parseIP(remoteAddr)
}

func parseIP(s string) (netip.Addr, bool) {
	ip, err := netip.ParseAddr(strings.TrimSpace(s))
	if err != nil {
		return netip.Addr{}, false
	}
	return ip.Unmap(), true
}

// requestID takes the caller's X-Request-ID when it is a sane token (so
// it can't flood or forge log lines), else mints a UUIDv7, and puts it on
// the context and the response.
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestIDHeader)
		if !validRequestID(id) {
			id = uuid.NewV7().String()
		}
		w.Header().Set(requestIDHeader, id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id)))
	})
}

func validRequestID(id string) bool {
	if id == "" || len(id) > 128 {
		return false
	}
	for _, c := range id {
		if c < '!' || c > '~' {
			return false
		}
	}
	return true
}

// routePattern is the matched chi route, "not_found" when none matched.
func routePattern(r *http.Request) string {
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		if p := rctx.RoutePattern(); p != "" {
			return p
		}
	}
	return "not_found"
}

func hasPrefix(path string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// statusWriter records the response status. It forwards Hijack (the SSE
// streams and websockets) and Unwrap, so http.ResponseController reaches
// the real writer to flush.
type statusWriter struct {
	http.ResponseWriter

	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Flush() {
	_ = http.NewResponseController(w.ResponseWriter).Flush()
}

func (w *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(w.ResponseWriter).Hijack()
}

func (w *statusWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }

// logs writes one http_request line per request outside skip.
func logs(logger *slog.Logger, skip []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if hasPrefix(r.URL.Path, skip) {
				next.ServeHTTP(w, r)
				return
			}
			start := time.Now()
			ww := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(ww, r)
			logger.Info("http_request",
				slog.String("request_id", RequestIDFromCtx(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", routePattern(r)),
				slog.Int("status", ww.status),
				slog.Duration("duration", time.Since(start)),
				slog.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

type httpCollectors struct {
	requests metric.Int64Counter
	duration metric.Float64Histogram
	problems metric.Int64Counter
}

// durationBuckets are 0.25ms doubling 16 times (~8s).
var durationBuckets = func() []float64 {
	b := make([]float64, 16)
	for i := range b {
		b[i] = 0.00025 * float64(int(1)<<i)
	}
	return b
}()

// newCollectors creates the HTTP instruments on meter. Names, attributes and
// buckets are what dashboards query: keep them stable. Names carry their
// Prometheus suffixes and no unit, so OTLP ingestion stores them as is.
func newCollectors(meter metric.Meter) (*httpCollectors, error) {
	requests, err := meter.Int64Counter("http_requests_total",
		metric.WithDescription("Total number of HTTP requests."))
	if err != nil {
		return nil, err
	}
	duration, err := meter.Float64Histogram("http_request_duration_seconds",
		metric.WithDescription("Histogram of HTTP response times."),
		metric.WithExplicitBucketBoundaries(durationBuckets...))
	if err != nil {
		return nil, err
	}
	problems, err := meter.Int64Counter("http_problems_total",
		metric.WithDescription("Problem answers by operation, problem type slug and status."))
	if err != nil {
		return nil, err
	}
	return &httpCollectors{requests: requests, duration: duration, problems: problems}, nil
}

// metrics counts and times requests per route outside skip. The status
// label is the status text ("OK"), as dashboards already expect. Realtime
// streams are counted but not timed (see realtime).
func metrics(c *httpCollectors, skip []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if hasPrefix(r.URL.Path, skip) {
				next.ServeHTTP(w, r)
				return
			}
			start := time.Now()
			ww := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(ww, r)
			route := routePattern(r)
			c.requests.Add(r.Context(), 1, metric.WithAttributes(
				attribute.String("route", route),
				attribute.String("method", r.Method),
				attribute.String("status", http.StatusText(ww.status)),
			))
			if !realtime(r) {
				c.duration.Record(r.Context(), time.Since(start).Seconds(),
					metric.WithAttributes(attribute.String("route", route)))
			}
		})
	}
}

// recoverer turns a handler panic into a logged 500.
func recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				rec := recover()
				if rec == nil {
					return
				}
				recErr, _ := rec.(error)
				if errors.Is(recErr, http.ErrAbortHandler) {
					panic(rec)
				}
				logger.Error("panic recovered",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Any("panic", rec),
					slog.String("stack", string(debug.Stack())),
					slog.String("remote_addr", r.RemoteAddr),
				)
				WriteError(w, r, fmt.Errorf("panic: %v", rec))
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// timeout bounds the request context; handlers stop at its deadline.
func timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func maxBodySize(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, limit)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ErrRateLimited is a client over its rate limit.
var ErrRateLimited = errors.New("too many requests from this client in a short time; wait a moment, then retry")

// limiterIdleTTL is how long a client's bucket outlives its last request.
// A bucket idle that long has refilled, so dropping it changes nothing.
const limiterIdleTTL = 3 * time.Minute

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// rateLimiter is a token bucket per client IP (r.RemoteAddr after realIP).
// Idle buckets are swept on the request path, so the map is bounded by the
// clients seen in the last limiterIdleTTL.
type rateLimiter struct {
	mu        sync.Mutex
	clients   map[string]*clientLimiter
	rps       rate.Limit
	burst     int
	lastSweep time.Time
	now       func() time.Time
}

func newRateLimiter(rl RateLimit) *rateLimiter {
	return &rateLimiter{
		clients: make(map[string]*clientLimiter),
		rps:     rate.Limit(rl.RPS),
		burst:   rl.Burst,
		now:     time.Now,
	}
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := rl.now()
	if now.Sub(rl.lastSweep) > limiterIdleTTL {
		for k, c := range rl.clients {
			if now.Sub(c.lastSeen) > limiterIdleTTL {
				delete(rl.clients, k)
			}
		}
		rl.lastSweep = now
	}
	c, ok := rl.clients[key]
	if !ok {
		c = &clientLimiter{limiter: rate.NewLimiter(rl.rps, rl.burst)}
		rl.clients[key] = c
	}
	c.lastSeen = now
	return c.limiter.AllowN(now, 1)
}

func (rl *rateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(r.RemoteAddr) {
			WriteError(w, r, ErrRateLimited)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CORSMaxAge is how long a browser may reuse a preflight answer: Chrome's
// cap. The edge Worker answers preflights with the same values
// (lib/ts/container-budget), so keep the two in step.
const CORSMaxAge = 2 * time.Hour

// CORS is the harness's CORS middleware for the given comma-separated
// origins and headers (empty headers: the standard set). Credentials are
// allowed.
func CORS(allowedOrigins, allowedHeaders string) func(http.Handler) http.Handler {
	headers := cleanCSV(allowedHeaders)
	if len(headers) == 0 {
		headers = []string{"Content-Type", "Authorization", requestIDHeader, "X-API-Key"}
	}
	return cors.New(cors.Options{
		AllowedOrigins:   cleanCSV(allowedOrigins),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   headers,
		AllowCredentials: true,
		MaxAge:           int(CORSMaxAge.Seconds()),
	}).Handler
}
