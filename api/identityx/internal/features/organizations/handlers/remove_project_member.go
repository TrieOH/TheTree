package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// RemoveProjectMember godoc
// @Summary Remove a project member
// @Description Lets you remove a member from the organization project
// @Tags organizations
// @ID organizations_removeprojectmember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RemoveOrgProjectMemberRequest true "Member details"
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{organization_id}/projects/{project_id}/members [delete]
func (h *Handlers) RemoveProjectMember(w http.ResponseWriter, r *http.Request) {
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
	var payload models.RemoveOrgProjectMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.RemoveProjectMember(r.Context(), payload.ToInput(orgID, projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
