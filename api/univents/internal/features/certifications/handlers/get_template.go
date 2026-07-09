package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetTemplateByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	tpl, err := h.queries.GetTemplateByID(r.Context(), templateID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tpl)
}
