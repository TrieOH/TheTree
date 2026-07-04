package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// List godoc
// @Summary List actors by project
// @Tags actors
// @ID actors_list
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.Actor}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /projects/{project_id}/actors [get]
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	actor, err := h.queries.List(r.Context(), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, actor)
}
