package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.DeleteTemplate(r.Context(), templateID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, nil, http.StatusNoContent)
}
