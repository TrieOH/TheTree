package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.UpdateCertificationTemplateRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	template, err := handler.commands.UpdateTemplate(r.Context(), payload.ToInput(templateID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, template)
}
