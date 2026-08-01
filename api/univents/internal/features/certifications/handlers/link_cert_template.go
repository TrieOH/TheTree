package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) LinkCertTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CertTemplateProgramRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.ops.LinkCertTemplate(r.Context(), templateID, payload.ProgramID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, nil, http.StatusCreated)
}
