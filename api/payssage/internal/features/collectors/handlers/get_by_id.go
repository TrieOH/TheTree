package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	collectorID, err := req.Path("collector_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	collector, err := h.queries.GetByID(r.Context(), collectorID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, collector, http.StatusOK)
}
