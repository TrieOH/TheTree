package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) ListCertTemplateLinks(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	links, err := handler.queries.ListCertTemplateLinks(r.Context(), templateID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, links)
}
