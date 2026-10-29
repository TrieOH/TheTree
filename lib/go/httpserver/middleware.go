package httpserver

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/cors"
	"go.uber.org/zap"
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
			id = uuid.Must(uuid.NewV7()).String()
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
func logs(logger *zap.Logger, skip []string) func(http.Handler) http.Handler {
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
				zap.String("request_id", RequestIDFromCtx(r.Context())),
				zap.String("method", r.Method),
				zap.String("path", routePattern(r)),
				zap.Int("status", ww.status),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}

type httpCollectors struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

// newCollectors registers the HTTP metrics. Names, labels and buckets are
// what dashboards query: keep them stable.
func newCollectors(reg prometheus.Registerer) (*httpCollectors, error) {
	c := &httpCollectors{
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		}, []string{"route", "method", "status"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP response times.",
			Buckets: prometheus.ExponentialBuckets(0.00025, 2, 16),
		}, []string{"route"}),
	}
	for _, col := range []prometheus.Collector{c.requests, c.duration} {
		err := reg.Register(col)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// metrics counts and times requests per route outside skip. The status
// label is the status text ("OK"), as dashboards already expect.
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
			c.requests.WithLabelValues(route, r.Method, http.StatusText(ww.status)).Inc()
			c.duration.WithLabelValues(route).Observe(time.Since(start).Seconds())
		})
	}
}

// recoverer turns a handler panic into a logged 500.
func recoverer(logger *zap.Logger) func(http.Handler) http.Handler {
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
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Any("panic", rec),
					zap.String("stack", string(debug.Stack())),
					zap.String("remote_addr", r.RemoteAddr),
				)
				writeMiddlewareError(w, r, http.StatusInternalServerError, fmt.Errorf("panic: %v", rec))
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

var errRateLimited = errors.New("rate limit exceeded")

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
			writeMiddlewareError(w, r, http.StatusTooManyRequests, errRateLimited)
			return
		}
		next.ServeHTTP(w, r)
	})
}

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
		MaxAge:           int((10 * time.Minute).Seconds()),
	}).Handler
}
