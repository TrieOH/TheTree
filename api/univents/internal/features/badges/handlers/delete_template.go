package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}

	err = h.commands.DeleteTemplate(r.Context(), templateID)
	if fun.Bail(w, err) {
		return
	}

	fun.NoContent().Send(w)
}
