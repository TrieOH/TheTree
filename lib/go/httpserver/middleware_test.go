package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type identityKey struct{}

func TestRealIPUsesCloudflareClientIP(t *testing.T) {
	cases := []struct {
		name, peer, cf, xff, want string
	}{
		{"CF-Connecting-IP wins", "10.0.0.2:40000", "198.51.100.1", "", "198.51.100.1"},
		{"X-Forwarded-For ignored", "10.0.0.2:40000", "", "198.51.100.2", "10.0.0.2"},
		{"no header: port stripped", "203.0.113.7:51234", "", "", "203.0.113.7"},
		{"garbage header: peer", "127.0.0.1:40000", "not-an-ip", "", "127.0.0.1"},
		{"ipv6 peer: port stripped", "[2001:db8::1]:443", "", "", "2001:db8::1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got string
			h := realIP(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { got = r.RemoteAddr }))
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
			req.RemoteAddr = tc.peer
			if tc.cf != "" {
				req.Header.Set("CF-Connecting-IP", tc.cf)
			}
			if tc.xff != "" {
				req.Header.Set("X-Forwarded-For", tc.xff)
			}
			h.ServeHTTP(httptest.NewRecorder(), req)
			if got != tc.want {
				t.Fatalf("RemoteAddr = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestRateLimitKeysPerClientNotConnection — a client reconnecting from a
// new port shares its bucket.
func TestRateLimitKeysPerClientNotConnection(t *testing.T) {
	h := testRouter(t, Config{RateLimit: RateLimit{RPS: 0.001, Burst: 1}})
	codes := make([]int, 0, 2)
	for _, port := range []string{"1111", "2222"} {
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil)
		req.RemoteAddr = "203.0.113.7:" + port
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		codes = append(codes, rec.Code)
	}
	if codes[0] != http.StatusOK || codes[1] != http.StatusTooManyRequests {
		t.Fatalf("codes = %v, want [200 429]", codes)
	}
}

func TestRateLimiterSweepsIdleClients(t *testing.T) {
	now := time.Unix(0, 0)
	rl := newRateLimiter(RateLimit{RPS: 1, Burst: 1})
	rl.now = func() time.Time { return now }

	rl.allow("a")
	rl.allow("b")
	now = now.Add(limiterIdleTTL + time.Second)
	rl.allow("c")

	if len(rl.clients) != 1 || rl.clients["c"] == nil {
		t.Fatalf("clients = %v, want only c", rl.clients)
	}
}

// TestLogsCarryGeneratedRequestID — a request without X-Request-ID is
// logged under the ID the harness minted and returned.
func TestLogsCarryGeneratedRequestID(t *testing.T) {
	core, logged := observer.New(zap.InfoLevel)
	h := requestID(logs(zap.New(core), nil)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/x", nil))

	id := rec.Header().Get(requestIDHeader)
	entries := logged.All()
	if id == "" || len(entries) != 1 || entries[0].ContextMap()["request_id"] != id {
		t.Fatalf("response id %q, log entries %v", id, entries)
	}
}

func TestRequestIDRejectsUnsafeValues(t *testing.T) {
	cases := map[string]bool{
		"abc-123":                    true,
		"":                           false,
		"has space":                  false,
		"line\nbreak":                false,
		strings.Repeat("a", 129):     false,
		"0190f6e4-7c1a-7000-8000-01": true,
	}
	for in, keep := range cases {
		h := requestID(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		req.Header[http.CanonicalHeaderKey(requestIDHeader)] = []string{in}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if got := rec.Header().Get(requestIDHeader); (got == in) != keep {
			t.Fatalf("X-Request-ID %q → %q, keep = %v", in, got, keep)
		}
	}
}

// TestStatusWriterFlushes — streaming handlers behind logs/metrics can
// still flush through http.ResponseController.
func TestStatusWriterFlushes(t *testing.T) {
	rec := httptest.NewRecorder()
	h := logs(zap.NewNop(), nil)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		err := http.NewResponseController(w).Flush()
		if err != nil {
			t.Fatalf("flush: %v", err)
		}
	}))
	h.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))
	if !rec.Flushed {
		t.Fatalf("recorder not flushed")
	}
}

func TestAuthenticatorsPrimitives(t *testing.T) {
	errBadKey := errors.New("invalid api key")
	p := Authenticators{
		JWT: func(ctx context.Context, token string) (context.Context, error) {
			return context.WithValue(ctx, identityKey{}, "jwt:"+token), nil
		},
		APIKey: func(ctx context.Context, key string) (context.Context, error) {
			if key != "good" {
				return ctx, errBadKey
			}
			return context.WithValue(ctx, identityKey{}, "key:"+key), nil
		},
	}.Primitives()

	cases := []struct {
		name    string
		mw      func(http.Handler) http.Handler
		headers map[string]string
		status  int
		who     string
	}{
		{"jwt ok", p.JWT, map[string]string{"Authorization": "Bearer t1"}, 200, "jwt:t1"},
		{"jwt missing", p.JWT, nil, 401, ""},
		{"jwt empty bearer", p.JWT, map[string]string{"Authorization": "Bearer "}, 401, ""},
		{"api key ok", p.APIKey, map[string]string{"X-API-Key": "good"}, 200, "key:good"},
		{"api key rejected", p.APIKey, map[string]string{"X-API-Key": "bad"}, 401, ""},
		{"any prefers api key", p.Any, map[string]string{"X-API-Key": "good", "Authorization": "Bearer t1"}, 200, "key:good"},
		{"any via jwt", p.Any, map[string]string{"Authorization": "Bearer t1"}, 200, "jwt:t1"},
		{"any with nothing", p.Any, nil, 401, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var who string
			h := tc.mw(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				who, _ = r.Context().Value(identityKey{}).(string)
			}))
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.status || who != tc.who {
				t.Fatalf("got (%d, %q), want (%d, %q)", rec.Code, who, tc.status, tc.who)
			}
		})
	}
}
