package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListByUser(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	userID, err := req.Path("user_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	certs, err := h.queries.ListByUser(r.Context(), userID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, certs)
}
