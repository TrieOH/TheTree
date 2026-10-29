package webhooks

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
)

// Capture holds the raw request as received, for webhook signature
// verification (HMAC covers the exact body bytes, before any decoding).
type Capture struct {
	Req  *http.Request
	Body []byte
}

type captureKey struct{}

// RawRequestMiddleware captures the raw body for webhook paths
// (/webhooks/...) before the generated strict handler decodes the body, and
// restores r.Body so decoding still works. Mirrors api/payssage's
// handler-side capture: the signature covers the envelope bytes Payssage
// actually POSTed, not the re-encoded typed struct.
func RawRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/webhooks/") {
			body, err := io.ReadAll(r.Body)
			if err == nil {
				r.Body = io.NopCloser(bytes.NewReader(body))
				ctx := context.WithValue(r.Context(), captureKey{}, &Capture{Req: r, Body: body})
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RawRequestFrom returns the captured raw request, or nil when the
// middleware did not run for this request.
func RawRequestFrom(ctx context.Context) *Capture {
	c, _ := ctx.Value(captureKey{}).(*Capture)
	return c
}
