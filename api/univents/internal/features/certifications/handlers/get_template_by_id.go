package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetTemplateByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	template, err := h.ops.GetTemplateByID(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, template)
}
