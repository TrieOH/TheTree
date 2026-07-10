package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// GetProfile godoc
// @Summary Get an actor's profile for a project
// @Tags profiles
// @ID profiles_getprofile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param actor_id path string true "Actor ID"
// @Success 200 {object} fun.Response{data=models.ActorProfile}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /projects/{project_id}/actors/{actor_id}/profile [get]
func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
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
	profile, err := h.queries.GetProfile(r.Context(), actorID, projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, profile)
}
