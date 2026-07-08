package handlers

import (
	"net/http"
	"univents/contracts"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) SetActivityTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	activityID, err := req.Path("activity_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload contracts.SetCertificationTemplateRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.SetActivityTemplate(r.Context(), activityID, payload.CertificationTemplateID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
