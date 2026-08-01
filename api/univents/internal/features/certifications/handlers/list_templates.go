package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListTemplates(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	templates, err := h.ops.ListTemplates(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, templates)
}
