package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) UpsertPlatformProfile(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	actorID, err := req.Path("actor_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.UpsertProfileRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	profile, err := h.commands.UpsertPlatformProfile(r.Context(), payload.ToInput(actorID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, profile)
}
