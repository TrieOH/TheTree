package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetTemplateByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	tpl, err := h.queries.GetTemplateByID(r.Context(), templateID, editionID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tpl)
}
