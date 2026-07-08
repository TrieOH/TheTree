package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListByTarget(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	targetType := req.Query("target_type").String()
	targetID, err := req.Query("target_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	certs, err := h.queries.ListByTarget(r.Context(), targetType, targetID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, certs)
}
