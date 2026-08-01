package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	var payload models.CreateOrganizationRequest
	if fun.BailInto(w, fun.From(r), &payload) {
		return
	}
	org, err := h.ops.Create(r.Context(), payload.ToInput())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, org, http.StatusCreated)
}
