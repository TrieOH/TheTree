package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// CreateProjectActor godoc
// @Summary Create actors in a project
// @Tags organizations
// @ID organizations_createprojectactor
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=models.Actor}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{org_id}/projects/{project_id}/actors [post]
func (h *Handlers) CreateProjectActor(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	orgID, err := req.Path("org_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateActorRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	actors, err := h.commands.CreateProjectActor(r.Context(), orgID, payload.ToInput(&projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, actors)
}
