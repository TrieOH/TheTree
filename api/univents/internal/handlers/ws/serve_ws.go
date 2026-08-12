package ws

import (
	"net/http"
)

// ServeWS is the raw `WS /ws?token=...` route (D9) — mounted by the app
// router outside the strict handler, because the handshake carries the
// one-time token in the query string (not an Authorization header). The
// token proves prior REST auth for this purchase; the service consumes it
// atomically (one-time), sends the snapshot frame, then relays live frames
// until a terminal event closes the socket.
func (h *Handlers) ServeWS(w http.ResponseWriter, r *http.Request) {
	h.ops.ServeWS(w, r)
}
