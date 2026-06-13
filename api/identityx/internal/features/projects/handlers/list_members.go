package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// ListMembers godoc
// @Summary Lists the project members
// @Description Gets a list of all the members added to the project, includes the owner
// @Tags projects
// @ID projects_listmembers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.ProjectMember}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /projects/{project_id}/members [get]
func (h *Handlers) ListMembers(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.queries.ListMembers(r.Context(), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
