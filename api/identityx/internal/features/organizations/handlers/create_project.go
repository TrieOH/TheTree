package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// CreateProject godoc
// @Summary Create an organization project
// @Description Creates a project scoped to the organization, with the org owner as the project owner
// @Tags organizations
// @ID organizations_createproject
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateOrgProjectRequest true "Project creation data"
// @Success 201 {object} fun.Response{data=models.Project}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{organization_id}/projects [post]
func (h *Handlers) CreateProject(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateOrgProjectRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	project, err := h.commands.CreateProject(r.Context(), payload.ToInput(orgID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, project, http.StatusCreated)
}
