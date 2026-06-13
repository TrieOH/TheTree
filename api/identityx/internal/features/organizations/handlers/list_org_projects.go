package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// ListProjects godoc
// @Summary Lists the organization projects
// @Description Gets a list of all the projects created in the organization
// @Tags organizations
// @ID organizations_listprojects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.Project}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{organization_id}/projects [get]
func (h *Handlers) ListProjects(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	projects, err := h.queries.ListOrgProjects(r.Context(), orgID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, projects)
}
