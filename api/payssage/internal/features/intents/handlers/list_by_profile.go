package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListByProfile(w http.ResponseWriter, r *http.Request) {
	intents, err := h.queries.ListByProfile(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, intents, http.StatusOK)
}
