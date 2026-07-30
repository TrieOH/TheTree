package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) UnlinkCertTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.UnlinkCertTemplate(r.Context(), templateID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, nil, http.StatusNoContent)
}
