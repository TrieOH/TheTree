package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// GetByID godoc
// @Summary Get actors by ID
// @Tags organizations
// @ID organizations_getactorbyid
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.Actor}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{organization_id}/projects/{project_id}/actors/{actor_id} [get]
func (h *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
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
	members, err := h.queries.GetByID(r.Context(), actorID, projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
