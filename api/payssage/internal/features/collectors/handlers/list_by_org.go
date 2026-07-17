package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListByOrg(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	collectors, err := h.queries.ListByOrg(r.Context(), orgID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, collectors, http.StatusOK)
}
