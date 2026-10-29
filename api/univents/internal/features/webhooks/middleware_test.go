package webhooks

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRawRequestMiddleware_CapturesAndRestoresBody pins the raw-body handoff
// the signature verification depends on: for /webhooks/ paths the exact
// body bytes are captured into the context (HMAC input) and r.Body is
// restored so the strict server can still decode the envelope.
func TestRawRequestMiddleware_CapturesAndRestoresBody(t *testing.T) {
	body := []byte(`{"intent_id":"1f6a0a6c-0d4b-4b5a-9f0e-7f1c3d2b4a5e","event_type":"payment.succeeded"}`)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/webhooks/payssage", bytes.NewReader(body))
	req.Header.Set("X-Payssage-Signature", "deadbeef")

	var captured *Capture
	var decoded []byte
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = RawRequestFrom(r.Context())
		decoded, _ = io.ReadAll(r.Body)
	})

	RawRequestMiddleware(next).ServeHTTP(httptest.NewRecorder(), req)

	if captured == nil {
		t.Fatal("capture missing for /webhooks/ path")
	}
	if !bytes.Equal(captured.Body, body) {
		t.Fatalf("captured body = %q, want %q", captured.Body, body)
	}
	if captured.Req.Header.Get("X-Payssage-Signature") != "deadbeef" {
		t.Fatalf("captured request lost the signature header")
	}
	if !bytes.Equal(decoded, body) {
		t.Fatalf("r.Body not restored for the strict decode: %q", decoded)
	}
}

// TestRawRequestMiddleware_SkipsOtherPaths ensures the middleware is a no-op
// outside /webhooks/ (it runs on every route, so it must not buffer bodies
// or stash captures elsewhere).
func TestRawRequestMiddleware_SkipsOtherPaths(t *testing.T) {
	body := []byte(`{"hello":"world"}`)
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/events", bytes.NewReader(body))

	var captured *Capture
	var decoded []byte
	next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		captured = RawRequestFrom(r.Context())
		decoded, _ = io.ReadAll(r.Body)
	})

	RawRequestMiddleware(next).ServeHTTP(httptest.NewRecorder(), req)

	if captured != nil {
		t.Fatalf("capture should be nil outside /webhooks/, got %+v", captured)
	}
	if !bytes.Equal(decoded, body) {
		t.Fatalf("body changed for non-webhook path: %q", decoded)
	}
}
