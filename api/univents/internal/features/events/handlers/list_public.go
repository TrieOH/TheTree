package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListPublic(w http.ResponseWriter, r *http.Request) {
	events, err := h.ops.ListPublic(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, events)
}
