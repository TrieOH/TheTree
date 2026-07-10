package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// UpsertProfile godoc
// @Summary Upsert an actor's profile for a project
// @Tags profiles
// @ID profiles_upsertprofile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param actor_id path string true "Actor ID"
// @Param body body models.UpsertProfileRequest true "Profile payload"
// @Success 200 {object} fun.Response{data=models.ActorProfile}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 422 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /projects/{project_id}/actors/{actor_id}/profile [put]
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
	profile, err := h.commands.UpsertProfile(r.Context(), payload.ToInput(actorID), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, profile)
}
