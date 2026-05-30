package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// ListProjectMembers godoc
// @Summary Lists the members of an organization project
// @Description Gets a list of all the members added to the project, includes inherited organization members
// @Tags organizations
// @ID organizations_listprojectmembers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.ProjectMember}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{organization_id}/projects/{project_id}/members [get]
func (h *Handlers) ListProjectMembers(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.queries.ListOrgProjectMembers(r.Context(), orgID, projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
