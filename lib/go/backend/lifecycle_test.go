package backend

import (
	"net/http"
	"testing"
	"time"
)

func TestNewServerTimeouts(t *testing.T) {
	srv := newServer(http.NewServeMux(), "8080")

	if srv.Addr != ":8080" {
		t.Fatalf("addr = %q, want :8080", srv.Addr)
	}
	if srv.ReadTimeout != 30*time.Second || srv.WriteTimeout != 60*time.Second || srv.IdleTimeout != 120*time.Second {
		t.Fatalf("timeouts = %v/%v/%v, want 30s/60s/120s", srv.ReadTimeout, srv.WriteTimeout, srv.IdleTimeout)
	}
}
