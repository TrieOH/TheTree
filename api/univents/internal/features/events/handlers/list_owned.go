package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListOwned(w http.ResponseWriter, r *http.Request) {
	events, err := h.ops.ListOwned(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, events)
}
