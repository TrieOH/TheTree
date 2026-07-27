package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	var payload models.CreateProjectRequest
	if bind.BailInto(w, fun.From(r), &payload) {
		return
	}
	org, err := h.commands.Create(r.Context(), payload.ToInput(nil))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, org, http.StatusCreated)
}
