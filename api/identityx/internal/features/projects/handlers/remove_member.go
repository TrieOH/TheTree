package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// RemoveMember godoc
// @Summary Remove a project member
// @Description Lets you remove a member from the project
// @Tags projects
// @ID projects_removemember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RemoveProjectMemberRequest true "Project details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /projects/{project_id}/members [delete]
func (h *Handlers) RemoveMember(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.RemoveProjectMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.RemoveMember(r.Context(), payload.ToInput(projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
