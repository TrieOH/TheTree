package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}

	template, err := h.ops.GetTemplate(r.Context(), templateID)
	if fun.Bail(w, err) {
		return
	}

	fun.Respond(w, template)
}
