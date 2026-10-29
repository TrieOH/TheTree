package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// mountRaw mounts the realtime routes (split 6) — deliberately outside the
// strict handler: the WS handshake cannot carry Authorization headers (the
// one-time query token is the auth), and SSE must stream without the
// strict server's validation machinery buffering the body. Both are documented in
// the spec's description block, not as spec ops.
func mountRaw(r chi.Router, h *server) {
	r.Get("/editions/{edition_id}/store/stream", h.ServeStoreStream)
	r.Handle("/ws", http.HandlerFunc(h.ServeWS))
}
