package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) SetEditionTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.SetCertificationTemplateRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.SetEditionTemplate(r.Context(), editionID, payload.CertificationTemplateID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
