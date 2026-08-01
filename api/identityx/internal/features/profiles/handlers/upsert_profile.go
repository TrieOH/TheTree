package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) UpsertProfile(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	actorID, err := req.Path("actor_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.UpsertProfileRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	profile, err := h.ops.UpsertProfile(r.Context(), payload.ToInput(actorID), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, profile)
}
